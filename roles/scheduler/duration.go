package scheduler

import (
	"fmt"
	"sync"
	"time"

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

	allAssignmentsExecutions := make(map[string][]time.Duration)
	var allDurations []time.Duration
	for _, a := range assignments {
		allAssignmentsExecutions[getAssignmentAsString(a, s.ExecutionType)] = allDurations
	}

	wg.Add(s.numberOfWorkers + 3)
	for i := 0; i < s.numberOfWorkers; i++ {
		go func(num int) {
			defer wg.Done()
			for a := range s.tasksChan {
				duration, err := s.executeTaskFromAssignment(&a)
				if err != nil {
					s.Logger.Error(fmt.Sprintf("worker could not run database query %v", &a), zap.Error(err))
				}
				title := getAssignmentAsString(a, s.ExecutionType)
				mutex.Lock()
				allAssignmentsExecutions = appendDurationToAssignmentResults(title, allAssignmentsExecutions, duration)
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
		Assignments: allAssignmentsExecutions,
		Durations:   concatAllDurations(allAssignmentsExecutions),
		// TODO: collect responses from api server by kind and total responses for each kind
		Response: nil,
		// TODO: collect error kind and total errors for each error kind
		Errors: nil,
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

func concatAllDurations(assignmentResults map[string][]time.Duration) []time.Duration {
	var allDurations []time.Duration
	for _, val := range assignmentResults {
		for _, dur := range val {
			allDurations = append(allDurations, dur)
		}
	}
	return allDurations
}
