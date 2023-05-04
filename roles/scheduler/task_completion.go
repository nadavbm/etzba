package scheduler

import (
	"fmt"
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

	var allDurations []time.Duration
	// Process results
	for r := range results {
		fmt.Println(r)
		allDurations = append(allDurations, r)
	}

	res := &Result{
		Assignments: assignments,
		Durations:   allDurations,
		// TODO: collect responses from api server by kind and total responses for each kind
		Response: nil,
		// TODO: collect error kind and total errors for each error kind
		Errors: nil,
	}

	return res, nil
}
