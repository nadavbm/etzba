package calculator

import (
	"time"

	"github.com/nadavbm/etzba/roles/scheduler"
)

// Calculator will calculate  and return float64 ms or time.Duration
type Calculator struct {
	Result scheduler.Result
}

// NewCalculator creates an instance of Calculator
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Durations compose of total # processed tasks, total processing time for the job, the minimum task time,
// the median  time, the average  time, and the maximum  time.
type Durations struct {
	Total                     int     `json:"total"`
	TotalJobTime              float64 `json:"job_time"`
	MinimumTime               float64 `json:"min_duration"`
	MedianTime                float64 `json:"med_duration"`
	AverageTime               float64 `json:"avg_duration"`
	MaximumTime               float64 `json:"max_duration"`
	TotalJobsOfAllWorkersTime float64
}

// GetResult get all required results in float64 as time.Duration is in type of float64
func (c *Calculator) GetResult(allDurations []time.Duration) *Durations {
	return &Durations{
		Total:        c.getTotalProcesessed(allDurations),
		TotalJobTime: c.getTotalProcesessedTime(allDurations),
		MinimumTime:  c.getMinimumTime(allDurations),
		MedianTime:   c.getMedianTime(allDurations),
		AverageTime:  c.getAverageTime(allDurations),
		MaximumTime:  c.getMaximumTime(allDurations),
	}
}

func (c *Calculator) getTotalProcesessed(allDurations []time.Duration) int {
	return len(allDurations)
}

func (c *Calculator) getTotalProcesessedTime(allDurations []time.Duration) float64 {
	var sum float64
	for _, dur := range allDurations {
		sum += float64(dur.Milliseconds())
	}
	return sum
}

func (c *Calculator) getMinimumTime(allDurations []time.Duration) float64 {
	var min int64
	for i, dur := range allDurations {
		if i == 0 {
			min = int64(dur.Milliseconds())
		}
		if min > int64(dur.Milliseconds()) {
			min = int64(dur.Milliseconds())
		}
	}
	return float64(min)
}

func (c *Calculator) getMedianTime(allDurations []time.Duration) float64 {
	total := c.getTotalProcesessed(allDurations)
	var median int
	if len(allDurations) < 2 {
		median = 0
	} else {
		median = int(total / 2)
	}
	return float64(allDurations[median].Milliseconds())
}

func (c *Calculator) getAverageTime(allDurations []time.Duration) float64 {
	total := c.getTotalProcesessedTime(allDurations)
	length := c.getTotalProcesessed(allDurations)
	return total / float64(length)
}

func (c *Calculator) getMaximumTime(allDurations []time.Duration) float64 {
	var max int64
	for _, dur := range allDurations {
		if max < int64(dur.Milliseconds()) {
			max = int64(dur.Milliseconds())
		}
	}
	return float64(max)
}
