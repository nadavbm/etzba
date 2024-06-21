package calculator

import (
	"time"
)

func GetTotalProcesessed(allDurations []time.Duration) int {
	return len(allDurations)
}

func GetTotalProcesessedTime(allDurations []time.Duration) float64 {
	var sum float64
	for _, dur := range allDurations {
		sum += float64(dur.Milliseconds())
	}
	return sum
}

func GetMinimumTime(allDurations []time.Duration) float64 {
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

func GetMedianTime(allDurations []time.Duration) float64 {
	total := GetTotalProcesessed(allDurations)
	var median int
	if len(allDurations) < 2 {
		median = 0
	} else {
		median = int(total / 2)
	}
	return float64(allDurations[median].Milliseconds())
}

func GetAverageTime(allDurations []time.Duration) float64 {
	total := GetTotalProcesessedTime(allDurations)
	length := GetTotalProcesessed(allDurations)
	return total / float64(length)
}

func GetMaximumTime(allDurations []time.Duration) float64 {
	var max int64
	for _, dur := range allDurations {
		if max < int64(dur.Milliseconds()) {
			max = int64(dur.Milliseconds())
		}
	}
	return float64(max)
}
