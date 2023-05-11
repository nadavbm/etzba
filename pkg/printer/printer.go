package printer

import (
	"fmt"

	"github.com/nadavbm/etzba/pkg/calculator"
	"github.com/nadavbm/etzba/roles/scheduler"
	"github.com/nadavbm/etzba/roles/worker"
)

func PrintToTerminal(r *scheduler.Result) {
	printAllTaskDurations(r)
	printDetailedAssignmentExecutions(r)
	printAllResponsesPerAssignment(r)
}

func printAllTaskDurations(r *scheduler.Result) {
	c := calculator.NewCalculator()
	result := c.GetResult(r.Durations)
	fmt.Println("\nresults durations: \n==================")
	fmt.Println("\t\ttotal processed tasks: \t", result.Total)
	fmt.Println(fmt.Sprintf("\t\ttotal processing time: \t%vms", result.TotalJobTime))
	fmt.Println(fmt.Sprintf("\t\tmin_duration: \t%vms", result.MinimumTime))
	fmt.Println(fmt.Sprintf("\t\tmed_duration: \t%vms", result.MedianTime))
	fmt.Println(fmt.Sprintf("\t\tavg_duration: \t%vms", result.AverageTime))
	fmt.Println(fmt.Sprintf("\t\tmax_duration: \t%vms", result.MaximumTime))
	return
}

func printDetailedAssignmentExecutions(r *scheduler.Result) {
	fmt.Println("\nassignment kinds during execution: \n==================================")
	assignmentMap := r.Assignments
	for a, t := range assignmentMap {
		c := calculator.NewCalculator()
		result := c.GetResult(t)
		fmt.Println(fmt.Sprintf("\t\t%s: \n\t\t%s", a, fmt.Sprintf("----------------------------------------------------------------")))
		fmt.Println(fmt.Sprintf("\t\ttotal executions: \t%d", result.Total))
		fmt.Println(fmt.Sprintf("\t\tavg_duration: \t%vms\n", result.AverageTime))
	}
}

func printAllResponsesPerAssignment(r *scheduler.Result) {
	fmt.Println("\nall service responses per assignment: \n==================================")
	assignmentAndResponsesMap := r.Responses
	for a, r := range assignmentAndResponsesMap {
		fmt.Println(fmt.Sprintf("\t\t%s: \n\t\t----------------------------------------------------------------", a))
		statusCount := getAllResponsesPerAssignment(r)
		for s, t := range statusCount {
			fmt.Println(fmt.Sprintf("\t\tstatus code: \t%d", s))
			fmt.Println(fmt.Sprintf("\t\ttotal: \t%d\n", t))
		}
	}
}

func getAllResponsesPerAssignment(responses []*worker.Response) map[int]int {
	statusCount := make(map[int]int)
	for _, r := range responses {
		_, exist := statusCount[r.Status]

		if exist {
			statusCount[r.Status] += 1
		} else {
			statusCount[r.Status] = 1
		}
	}
	return statusCount
}
