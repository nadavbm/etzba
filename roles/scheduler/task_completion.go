package scheduler

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

func (s *Scheduler) ExecuteJobUntilCompletion() (*Result, error) {
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

	results := make(chan time.Duration)
	workCh := make(workerChannel)

	// Start workers
	var wg sync.WaitGroup
	wg.Add(s.numberOfWorkers)
	for i := 0; i < s.numberOfWorkers; i++ {
		go func(num int) {
			defer wg.Done()
			for a := range workCh {
				duration, err := s.executeTaskFromAssignment(&a)
				if err != nil {
					s.Logger.Fatal("could not execute task from assignment", zap.Error(err))
				}
				results <- duration

				title := getAssignmentAsString(a, s.ExecutionType)
				mutex.Lock()
				allAssignmentsExecutions = appendDurationToAssignmentResults(title, allAssignmentsExecutions, duration)
				mutex.Unlock()
			}
		}(i)
	}

	// Close result channel when workers done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Send work to be done
	go func() {
		for _, a := range assignments {
			workCh <- a
		}
		close(workCh)
	}()

	// Process results
	for r := range results {
		allDurations = append(allDurations, r)
	}

	res := &Result{
		//Assignments: map[getAssignmentStrin assignments,
		Assignments: allAssignmentsExecutions,
		Durations:   allDurations,
		// TODO: collect responses from api server by kind and total responses for each kind
		Response: nil,
		// TODO: collect error kind and total errors for each error kind
		Errors: nil,
	}

	return res, nil
}
