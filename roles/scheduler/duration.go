package scheduler

import (
	"fmt"
	"time"

	"github.com/nadavbm/etzba/roles/apiclient"
	"github.com/nadavbm/etzba/roles/common"
	"github.com/nadavbm/etzba/roles/prompusher"
	"github.com/nadavbm/etzba/roles/worker"
	"go.uber.org/zap"
)

const sleepTimeBeforeClosingChannels = 3

// ExecuteJobByDuration when "--duration=Xx" is given via command line, shceduler will fill work channel with assignments until the job duration is over
// after execution is completed, it will return the result of the load test
func (s *Scheduler) ExecuteJobByDuration() (*common.Result, error) {
	assignments, err := s.setAssignmentsToWorkers()
	if err != nil {
		s.Logger.Fatal("could not create assignments", zap.Error(err))
		panic(err)
	}

	allAssignmentsExecutionsDurations, allAssignmentsExecutionsResponses, err := s.prepareAssignmentsForResultCollection(assignments)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	wg.Add(s.Settings.NumberOfWorkers + 3)
	for i := 0; i < s.Settings.NumberOfWorkers; i++ {
		go func(num int) {
			defer wg.Done()
			for a := range s.tasksChan {
				if s.Settings.Verbose {
					s.Logger.Info(fmt.Sprintf("worker %d execute task %v", num, &a))
				}
				duration, resp, err := s.executeTaskFromAssignment(&a)
				if err != nil {
					s.Logger.Error(fmt.Sprintf("worker could not execute task %v", &a), zap.Error(err))
				}
				title := getAssignmentAsString(a, s.Settings.ExecutionType)
				mutex.Lock()
				allAssignmentsExecutionsDurations = appendDurationsToAssignmentResults(title, allAssignmentsExecutionsDurations, duration)
				allAssignmentsExecutionsResponses = appendResponsesToAssignmentResults(title, allAssignmentsExecutionsResponses, resp)
				mutex.Unlock()
			}
		}(i)
	}

	go s.addToWorkChannel(s.Settings.SetRps(), s.Settings.Duration, s.tasksChan, assignments)

	go func() {
		wg.Wait()
	}()

	for {
		val, ok := <-s.tasksChan
		if !ok {
			wg.Done()
			break
		} else {
			s.tasksChan <- val
		}
	}

	elapsed := time.Since(now) - time.Duration(sleepTimeBeforeClosingChannels*(time.Second))
	s.Logger.Info("Calculating results", zap.Any("elapsed", elapsed.Seconds()))
	pushJobDurationToPrometheus(nil, elapsed.Seconds())
	return common.PrepareResultOuput("", s.Settings.ExecutionType, elapsed, allAssignmentsExecutionsDurations, allAssignmentsExecutionsResponses), nil
}

// addToWorkChannel will add assignments to work channel and close the channel when the duration time is over
func (s *Scheduler) addToWorkChannel(sleepTime, duration time.Duration, c chan worker.Assignment, assigments []worker.Assignment) {
	defer wg.Done()
	timer := time.NewTimer(duration)

	for {
		select {
		case <-timer.C:
			timer.Stop()
			s.Logger.Info("job completed after", zap.Any("duration", duration.Seconds()))
			wg.Done()
			// TODO: set max query time and sleep before closing the channl to allow all workers finish their assignment executions.
			time.Sleep(sleepTimeBeforeClosingChannels * time.Second)
			close(c)
			return
		default:
			for _, a := range assigments {
				time.Sleep(sleepTime)
				c <- a
			}
		}
	}
}

// appendDurationsToAssignmentResults collect all durations per assignment during job execution for any worker and return a map with all assignments and their durations
func appendDurationsToAssignmentResults(title string, assignmentResults map[string][]time.Duration, duration time.Duration) map[string][]time.Duration {
	for key, val := range assignmentResults {
		if title == key {
			val = append(val, duration)
			assignmentResults[key] = val
		}
	}

	return assignmentResults
}

// appendResponsesToAssignmentResults collect all responses from server per assignment
func appendResponsesToAssignmentResults(title string, assignmentResponses map[string][]*apiclient.Response, response *apiclient.Response) map[string][]*apiclient.Response {
	for key, val := range assignmentResponses {
		if title == key {
			val = append(val, response)
			assignmentResponses[key] = val
		}
	}

	return assignmentResponses
}

func pushJobDurationToPrometheus(labels []string, value float64) {
	jobDuration := prompusher.PrometheusClient.NewHistogram("job_duration", "count total executions from result", labels)
	prompusher.PrometheusClient.PushHistogram(jobDuration, "etzba_result", "job_duration", labels, value)
}
