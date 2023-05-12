package scheduler

import (
	"fmt"
	"sync"
	"time"

	"github.com/nadavbm/etzba/roles/apiclient"
	"github.com/nadavbm/etzba/roles/worker"
	"go.uber.org/zap"
)

var wg sync.WaitGroup
var mutex = &sync.Mutex{}

func (s *Scheduler) ExecuteTaskByDuration() (*Result, error) {
	assignments, err := s.setAssignmentsToWorkers()
	if err != nil {
		s.Logger.Fatal("could not create assignments")
		return nil, err
	}

	allAssignmentsExecutionsDurations := make(map[string][]time.Duration)
	allAssignmentsExecutionsResponses := make(map[string][]*apiclient.Response)
	var allAPIResponses []*apiclient.Response
	var allDurations []time.Duration
	for _, a := range assignments {
		allAssignmentsExecutionsDurations[getAssignmentAsString(a, s.ExecutionType)] = allDurations
		allAssignmentsExecutionsResponses[getAssignmentAsString(a, s.ExecutionType)] = allAPIResponses
	}

	wg.Add(s.numberOfWorkers + 3)
	for i := 0; i < s.numberOfWorkers; i++ {
		go func(num int) {
			defer wg.Done()
			for a := range s.tasksChan {
				duration, resp, err := s.executeTaskFromAssignment(&a)
				if err != nil {
					s.Logger.Error(fmt.Sprintf("worker could not execute task %v", &a), zap.Error(err))
				}
				title := getAssignmentAsString(a, s.ExecutionType)
				mutex.Lock()
				allAssignmentsExecutionsDurations = appendDurationToAssignmentResults(title, allAssignmentsExecutionsDurations, duration)
				allAssignmentsExecutionsResponses = appendResponseToAssignmentResults(title, allAssignmentsExecutionsResponses, resp)
				mutex.Unlock()
			}
		}(i)
	}

	go addToWorkChannel(s.jobDuration, s.tasksChan, assignments)

	go func() {
		wg.Wait()
	}()

	for {
		val, ok := <-s.tasksChan
		if ok == false {
			wg.Done()
			break
		} else {
			s.tasksChan <- val
		}
	}

	res := &Result{
		Assignments: allAssignmentsExecutionsDurations,
		Durations:   concatAllDurations(allAssignmentsExecutionsDurations),
		Responses:   allAssignmentsExecutionsResponses,
	}

	return res, nil

}

func addToWorkChannel(duration time.Duration, c chan worker.Assignment, assigments []worker.Assignment) {
	defer wg.Done()
	timer := time.NewTimer(duration)

	for {
		select {
		case <-timer.C:
			timer.Stop()
			fmt.Println("time is over")
			wg.Done()
			time.Sleep(1 * time.Second)
			close(c)
			return
		default:
			for _, a := range assigments {
				c <- a
			}
		}
	}
}

func appendDurationToAssignmentResults(title string, assignmentResults map[string][]time.Duration, duration time.Duration) map[string][]time.Duration {
	for key, val := range assignmentResults {
		if title == key {
			val = append(val, duration)
			assignmentResults[key] = val
		}
	}

	return assignmentResults
}

func appendResponseToAssignmentResults(title string, assignmentResponses map[string][]*apiclient.Response, response *apiclient.Response) map[string][]*apiclient.Response {
	for key, val := range assignmentResponses {
		if title == key {
			val = append(val, response)
			assignmentResponses[key] = val
		}
	}

	return assignmentResponses
}

func concatAllDurations(assignmentResults map[string][]time.Duration) []time.Duration {
	var allDurations []time.Duration
	for _, val := range assignmentResults {
		for _, dur := range val {
			allDurations = append(allDurations, dur)
		}
	}
	return allDurations
}
