package scheduler

import (
	"fmt"
	"sync"
	"time"

	"github.com/nadavbm/etzba/pkg/debug"
	"github.com/nadavbm/etzba/roles/worker"
)

func (s *Scheduler) ExecuteJobUntilCompletion() error {
	data, err := worker.ReadCSVFile(s.HelperFile)
	if err != nil {
		s.Logger.Fatal("could not read csv file")
		return err
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
					s.Logger.Fatal("could not create worker")
				}
				debug.Debug("assignment and worker", a, num)
				duration, err := worker.GetSQLQueryDuration(&a)
				if err != nil {
					debug.Debug(err)
					s.Logger.Fatal("could not get sql query duration")
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
	fmt.Println("all durations:", allDurations)
	return nil
}
