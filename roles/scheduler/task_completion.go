package scheduler

import (
	"sync"
	"time"

	"github.com/nadavbm/etzba/roles/worker"
	"go.uber.org/zap"
)

func (s *Scheduler) ExecuteJobUntilCompletion() (*Result, error) {
	data, err := worker.ReadCSVFile(s.HelperFile)
	if err != nil {
		s.Logger.Fatal("could not read csv file")
		return nil, err
	}

	assignments := worker.SetSQLAssignmentsToWorkers(data)

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
				worker, err := worker.NewSQLWorker(s.Logger, s.ConfigFile)
				if err != nil {
					s.Logger.Fatal("could not create worker", zap.Error(err))
				}
				duration, err := worker.GetSQLQueryDuration(&a)
				if err != nil {
					s.Logger.Fatal("could not get sql query duration", zap.Error(err))
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
