package scheduler

import (
	"time"

	"github.com/nadavbm/etzba/pkg/debug"
	"github.com/nadavbm/etzba/roles/worker"
)

func (s *Scheduler) ExecuteTaskByDurationBak() ([]time.Duration, error) {
	data, err := worker.ReadCSVFile(s.HelperFile)
	if err != nil {
		s.Logger.Fatal("could not read csv file")
		return nil, err
	}

	assignments := worker.SetSQLAssignmentsToWorkers(data)

	durationEnd := make(chan bool)

	go func() {
		for {
			select {
			case <-durationEnd:
				debug.Debug("duration end")
				return
			default:
				s.tasksChan = s.setWorkChannelForDuration(assignments)
			}
		}
	}()

	time.Sleep(s.jobDuration)
	durationEnd <- true

	// receive all durations of all queries from the results channel

	for a := range s.tasksChan {
		assignments = append(assignments, a)
	}

	debug.Debug(assignments)
	// we sleep one second to allow all channels (fan out and fan in) to close properly and avoid panic
	time.Sleep(time.Second * 1)
	return nil, nil
}

func (s *Scheduler) setWorkChannelForDuration(assignments []worker.Assignment) <-chan worker.Assignment {
	debug.Debug("start duration", s)
	workCh := make(chan worker.Assignment)
	done := make(chan bool)
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				debug.Debug(" tick print", t)
				for _, a := range assignments {
					debug.Debug("assignment to channel", a)
					workCh <- a
				}
			}
		}
	}()

	debug.Debug("duration", s.jobDuration)
	time.Sleep(10 * time.Second)
	ticker.Stop()
	debug.Debug("end of job duration")
	done <- true
	debug.Debug("return")

	return workCh
}
