package scheduler

import (
	"fmt"
	"sync"
	"time"

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

	var allDurations []time.Duration

	wg.Add(s.numberOfWorkers + 3)
	for i := 0; i < s.numberOfWorkers; i++ {
		go func(num int) {
			defer wg.Done()
			for a := range s.tasksChan {
				duration, err := s.executeTaskFromAssignment(&a)
				if err != nil {
					s.Logger.Error(fmt.Sprintf("worker could not run database query %v", &a), zap.Error(err))
				}
				allDurations = append(allDurations, duration)
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
			fmt.Println(val, ok, "loop break")
			break
		} else {
			s.tasksChan <- val
		}
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
