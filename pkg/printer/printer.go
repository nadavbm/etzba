package printer

import (
	"fmt"

	"github.com/nadavbm/etzba/roles/calculator"
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
	kinds := filterDuplicateAssignments(r.Assignments, command)
	for i, k := range kinds {
		switch {
		case command == "api":
			fmt.Println(fmt.Sprintf("\t\t%d. %s", i, k))
		case command == "sql":
			fmt.Println(fmt.Sprintf("\t\t%d. %s", i, k))
			fmt.Println(fmt.Sprintf("\t\t  total executions %v", len(r.Assignments)))
		}
	}
}

func filterDuplicateAssignments(assignments []worker.Assignment, command string) []string {
	var filteredAssignments []string

	allAssignmentKind := make(map[string]bool)
	for _, a := range assignments {
		if _, value := allAssignmentKind[getAssignmentAsString(a, command)]; !value {
			allAssignmentKind[getAssignmentAsString(a, command)] = true
			filteredAssignments = append(filteredAssignments, getAssignmentAsString(a, command))
		}
	}
	return filteredAssignments
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

//func getTotalExecutionsPerAssignmentKind(assignments []worker.Assignment, command, aStr string) map[string]int {
//	count = 0
//	switch {
//	case command == "api":
//		return fmt.Sprintf("%v", a.ApiRequest)
//	case command == "sql":
//		return fmt.Sprintf("%s", sqlclient.ToSQL(&a.SqlQuery))
//	}
//	return 0
//}
