package printer

import (
	"fmt"

	"github.com/nadavbm/etzba/roles/calculator"
	"github.com/nadavbm/etzba/roles/scheduler"
)

func PrintTaskDurations(r *scheduler.Result) {
	c := calculator.NewCalculator()
	result := c.GetResult(r.Durations)
	fmt.Println("\n\nTasks processes in total: \t", result.Total)
	fmt.Println(fmt.Sprintf("\nTotal processing time for all queries in milliseconds: \t%vms", result.TotalOperationTime))
	fmt.Println(fmt.Sprintf("\nMinimum query time in milliseconds: \t%vms", result.MinimumTime))
	fmt.Println(fmt.Sprintf("\nMedian query time in milliseconds: \t%vms", result.MedianTime))
	fmt.Println(fmt.Sprintf("\nAverage query time in milliseconds: \t%vms", result.AverageTime))
	fmt.Println(fmt.Sprintf("\nMaximum query time in milliseconds: \t%vms", result.MaximumTime))

	// JobTimeEnd := time.Since(jobTimeStart)
	// result.TotalJobsOfAllWorkersTime = float64(JobTimeEnd.Milliseconds())
	// fmt.Println(fmt.Sprintf("\nThe job took: \t%vms", r.TotalJobsOfAllWorkersTime))
	// fmt.Println(fmt.Sprintf("length of all assignments %d, all calculated durations %d", len(assignments), len(allDurations)))
	return
}
