package scheduler

import (
	"fmt"
	"sync"
	"time"

	"github.com/nadavbm/etzba/pkg/debug"
	"github.com/nadavbm/etzba/roles/worker"
)

var wg sync.WaitGroup

func (s *Scheduler) ExecuteTaskByDuration() ([]time.Duration, error) {
	data, err := worker.ReadCSVFile(s.HelperFile)
	if err != nil {
		s.Logger.Fatal("could not read csv file")
		return nil, err
	}

	assignments := worker.SetSQLAssignmentsToWorkers(data)

	c := make(chan worker.Assignment, s.numberOfWorkers)
	var allDurations []time.Duration

	wg.Add(s.numberOfWorkers + 3)
	for i := 0; i < s.numberOfWorkers; i++ {
		go func(num int) {
			defer wg.Done()
			for w := range c {
				fmt.Println("assignment ", w, " worker ", num)
				worker, err := worker.NewSQLWorker(s.Logger, s.ConfigFile, s.HelperFile)
				if err != nil {
					s.Logger.Fatal("could not create worker")
				}
				duration, err := worker.GetSQLQueryDuration(&w)
				if err != nil {
					s.Logger.Error(fmt.Sprintf("worker could not run database query %v", &w))
				}

				debug.Debug("duration", duration)

				allDurations = append(allDurations, duration)
			}
		}(i)
	}

	go addToWorkChannel(s.jobDuration, c, assignments)

	go func() {
		wg.Wait()
	}()

	for {
		val, ok := <-c
		if ok == false {
			wg.Done()
			fmt.Println(val, ok, "loop break")
			break
		} else {
			c <- val
		}
	}
	return allDurations, nil

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
