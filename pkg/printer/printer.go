package printer

import (
	"fmt"

	"github.com/nadavbm/etzba/roles/calculator"
	"github.com/nadavbm/etzba/roles/scheduler"
)

func PrintTaskDurations(r *scheduler.Result) {
	c := calculator.NewCalculator()
	result := c.GetResult(r.Durations)
	fmt.Println("\nresults durations: \n==================")
	fmt.Println("\t\ttotal processed tasks: \t", result.Total)
	fmt.Println(fmt.Sprintf("\t\ttotal processing time: \t%vms", result.TotalOperationTime))
	fmt.Println(fmt.Sprintf("\t\tmin: \t%vms", result.MinimumTime))
	fmt.Println(fmt.Sprintf("\t\tmed: \t%vms", result.MedianTime))
	fmt.Println(fmt.Sprintf("\t\tavg: \t%vms", result.AverageTime))
	fmt.Println(fmt.Sprintf("\t\tmax: \t%vms", result.MaximumTime))

	// JobTimeEnd := time.Since(jobTimeStart)
	// result.TotalJobsOfAllWorkersTime = float64(JobTimeEnd.Milliseconds())
	// fmt.Println(fmt.Sprintf("\nThe job took: \t%vms", r.TotalJobsOfAllWorkersTime))
	// fmt.Println(fmt.Sprintf("length of all assignments %d, all calculated durations %d", len(assignments), len(allDurations)))
	return
}

func PrintAssignmentsKind(r *scheduler.Result, command string) {
	fmt.Println("\nassignment kinds during execution: \n==================================")
	for i, a := range r.Assignments {
		switch {
		case command == "api":
			fmt.Println(fmt.Sprintf("\t\t%d. %v", i, a.ApiRequest))
		case command == "sql":
			fmt.Println(fmt.Sprintf("\t\t%d. %v", i, a.SqlQuery))
		}
	}
}
