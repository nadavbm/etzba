package scheduler

import (
	"fmt"
	"time"

	"github.com/nadavbm/etzba/roles/common"
	"go.uber.org/zap"
)

// ExecuteJobUntilCompletion when omitting '--duration' from the command, this function will execute
// all assignments from the config file until all assignments completed
func (s *Scheduler) ExecuteJobUntilCompletion() (*common.Result, error) {
	assignments, err := s.setAssignmentsToWorkers()
	if err != nil {
		s.Logger.Fatal("could not create assignments", zap.Error(err))
		panic(err)
	}

	allAssignmentsExecutionsDurations, allAssignmentsExecutionsResponses, err := s.prepareAssignmentsForResultCollection(assignments)
	if err != nil {
		panic(err)
	}

	results := make(chan time.Duration)
	workCh := make(workerChannel)
	now := time.Now()

	// Start workers

	wg.Add(s.Settings.NumberOfWorkers)
	for i := 0; i < s.Settings.NumberOfWorkers; i++ {
		go func(num int) {
			defer wg.Done()
			for a := range workCh {
				if s.Settings.Verbose {
					s.Logger.Info(fmt.Sprintf("worker %d execute task %v", num, &a))
				}

				duration, resp, err := s.executeTaskFromAssignment(&a)
				if err != nil {
					s.Logger.Fatal("could not execute task from assignment", zap.Error(err))
				}

				title := getAssignmentAsString(a, s.Settings.ExecutionType)
				mutex.Lock()
				allAssignmentsExecutionsDurations = appendDurationsToAssignmentResults(title, allAssignmentsExecutionsDurations, duration)
				allAssignmentsExecutionsResponses = appendResponsesToAssignmentResults(title, allAssignmentsExecutionsResponses, resp)
				mutex.Unlock()

				results <- duration
			}
		}(i)
	}

	// Close result channel when workers done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Send work to be done
	rpsSleepTime := s.Settings.SetRps()
	go func() {
		for _, a := range assignments {
			time.Sleep(rpsSleepTime)
			workCh <- a
		}
		close(workCh)
	}()

	var allDurations []time.Duration
	// Process results before return (This is needed to complete all the tasks and close channels gracefully)
	for r := range results {
		allDurations = append(allDurations, r)
	}

	return common.PrepareResultOuput("", s.Settings.ExecutionType, time.Since(now), allAssignmentsExecutionsDurations, allAssignmentsExecutionsResponses), nil
}
