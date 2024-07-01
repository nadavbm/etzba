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
		sum += float64(dur.Seconds() * 1000)
	}
	return sum
}

func GetMinimumTime(allDurations []time.Duration) float64 {
	min := float64(allDurations[0].Seconds() * 1000)
	for i := 0; i < len(allDurations); i++ {
		if min > float64(allDurations[i].Seconds()*1000) {
			min = float64(allDurations[i].Seconds() * 1000)
		}
	}
	return min
}

func GetMedianTime(allDurations []time.Duration) float64 {
	total := GetTotalProcesessed(allDurations)
	var median int
	if len(allDurations) < 2 {
		median = 0
	} else {
		median = int(total / 2)
	}
	return float64(allDurations[median].Seconds() * 1000)
}

func GetAverageTime(allDurations []time.Duration) float64 {
	total := GetTotalProcesessedTime(allDurations)
	length := GetTotalProcesessed(allDurations)
	return total / float64(length)
}

func GetMaximumTime(allDurations []time.Duration) float64 {
	var max float64
	for _, dur := range allDurations {
		if max < float64(dur.Seconds()*1000) {
			max = float64(dur.Seconds() * 1000)
		}
	}
	return max
}
