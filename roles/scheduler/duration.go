package scheduler

import (
	"fmt"
	"sync"
	"time"

	"github.com/nadavbm/etzba/pkg/debug"
	"github.com/nadavbm/etzba/roles/worker"
	"go.uber.org/zap"
)

var wg sync.WaitGroup

func (s *Scheduler) ExecuteTaskByDuration() (*Result, error) {
	data, err := worker.ReadCSVFile(s.HelperFile)
	if err != nil {
		s.Logger.Fatal("could not read csv file")
		return nil, err
	}

	assignments := worker.SetSQLAssignmentsToWorkers(data)

	workCh := make(chan worker.Assignment, s.numberOfWorkers)
	var allDurations []time.Duration

	wg.Add(s.numberOfWorkers + 3)
	for i := 0; i < s.numberOfWorkers; i++ {
		go func(num int) {
			defer wg.Done()
			for a := range workCh {
				fmt.Println("assignment ", a, " worker ", num)
				duration, err := s.executeSQLQueriesFromAssignment(&a)
				if err != nil {
					s.Logger.Error(fmt.Sprintf("worker could not run database query %v", &a), zap.Error(err))
				}

				debug.Debug("duration", duration)

				allDurations = append(allDurations, duration)
			}
		}(i)
	}

	go addToWorkChannel(s.jobDuration, workCh, assignments)

	go func() {
		wg.Wait()
	}()

	for {
		val, ok := <-workCh
		if ok == false {
			wg.Done()
			fmt.Println(val, ok, "loop break")
			break
		} else {
			workCh <- val
		}
	}

	res := &Result{
		Durations: allDurations,
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

func (s *Scheduler) executeSQLQueriesFromAssignment(assignment *worker.Assignment) (time.Duration, error) {
	worker, err := worker.NewSQLWorker(s.Logger, s.ConfigFile)
	if err != nil {
		s.Logger.Fatal("could not create worker")
	}
	return worker.GetSQLQueryDuration(assignment)
}
