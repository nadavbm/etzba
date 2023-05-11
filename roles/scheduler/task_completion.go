package scheduler

import (
	"sync"
	"time"

	"github.com/nadavbm/etzba/roles/apiclient"
	"go.uber.org/zap"
)

func (s *Scheduler) ExecuteJobUntilCompletion() (*Result, error) {
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

	results := make(chan time.Duration)
	workCh := make(workerChannel)

	// Start workers
	var wg sync.WaitGroup
	wg.Add(s.numberOfWorkers)
	for i := 0; i < s.numberOfWorkers; i++ {
		go func(num int) {
			defer wg.Done()
			for a := range workCh {
				duration, resp, err := s.executeTaskFromAssignment(&a)
				if err != nil {
					s.Logger.Fatal("could not execute task from assignment", zap.Error(err))
				}
				results <- duration

				title := getAssignmentAsString(a, s.ExecutionType)
				mutex.Lock()
				allAssignmentsExecutionsDurations = appendDurationToAssignmentResults(title, allAssignmentsExecutionsDurations, duration)
				allAssignmentsExecutionsResponses = appendResponseToAssignmentResults(title, allAssignmentsExecutionsResponses, resp)
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
		Assignments: allAssignmentsExecutionsDurations,
		Durations:   concatAllDurations(allAssignmentsExecutionsDurations),
		Responses:   allAssignmentsExecutionsResponses,
	}

	return res, nil
}
