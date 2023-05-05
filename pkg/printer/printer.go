package printer

import (
	"fmt"

	"github.com/nadavbm/etzba/pkg/calculator"
	"github.com/nadavbm/etzba/roles/scheduler"
	"github.com/nadavbm/etzba/roles/sqlclient"
	"github.com/nadavbm/etzba/roles/worker"
)

func PrintTaskDurations(r *scheduler.Result) {
	c := calculator.NewCalculator()
	result := c.GetResult(r.Durations)
	fmt.Println("\nresults durations: \n==================")
	fmt.Println("\t\ttotal processed tasks: \t", result.Total)
	fmt.Println(fmt.Sprintf("\t\ttotal processing time: \t%vms", result.TotalOperationTime))
	fmt.Println(fmt.Sprintf("\t\tmin_duration: \t%vms", result.MinimumTime))
	fmt.Println(fmt.Sprintf("\t\tmed_duration: \t%vms", result.MedianTime))
	fmt.Println(fmt.Sprintf("\t\tavg_duration: \t%vms", result.AverageTime))
	fmt.Println(fmt.Sprintf("\t\tmax_duration: \t%vms", result.MaximumTime))
	return
}

func PrintAssignmentsKind(r *scheduler.Result, command string) {
	fmt.Println("\nassignment kinds during execution: \n==================================")
	assignmentMap := r.Assignments
	for a, t := range assignmentMap {
		c := calculator.NewCalculator()
		result := c.GetResult(t)
		fmt.Println(fmt.Sprintf("\t\t%s: \n\t\t%s", a, fmt.Sprintf("----------------------------------------------------------------")))
		fmt.Println(fmt.Sprintf("\t\ttotal executions: \t%d", result.Total))
		fmt.Println(fmt.Sprintf("\t\tavg_duration: \t%v\n", result.AverageTime))
	}
}

func getAssignmentAsString(a worker.Assignment, command string) string {
	switch {
	case command == "api":
		return fmt.Sprintf("%v", a.ApiRequest)
	case command == "sql":
		return fmt.Sprintf("%s", sqlclient.ToSQL(&a.SqlQuery))
	}
	return ""
}
