package printer

import (
	"fmt"

	"github.com/nadavbm/etzba/pkg/calculator"
	"github.com/nadavbm/etzba/roles/apiclient"
	"github.com/nadavbm/etzba/roles/scheduler"
)

func PrintToTerminal(r *scheduler.Result, collect bool) {
	printAllTaskDurations(r)
	printDetailedAssignmentExecutions(r)
	if collect {
		printAllResponsesPerAssignment(r)
	}
}

func printAllTaskDurations(r *scheduler.Result) {
	c := calculator.NewCalculator()
	result := c.GetResult(r.Durations)
	fmt.Println("\nresults durations: \n==================")
	fmt.Println("\ttotal processed tasks: \t", result.Total)
	fmt.Println(fmt.Sprintf("\ttotal processing time: \t%vms", result.TotalJobTime))
	fmt.Println(fmt.Sprintf("\tmin_duration: \t%vms", result.MinimumTime))
	fmt.Println(fmt.Sprintf("\tmed_duration: \t%vms", result.MedianTime))
	fmt.Println(fmt.Sprintf("\tavg_duration: \t%vms", result.AverageTime))
	fmt.Println(fmt.Sprintf("\tmax_duration: \t%vms", result.MaximumTime))
	return
}

func printDetailedAssignmentExecutions(r *scheduler.Result) {
	fmt.Println("\nassignment kinds during execution: \n==================================")
	assignmentMap := r.Assignments
	for a, t := range assignmentMap {
		c := calculator.NewCalculator()
		result := c.GetResult(t)
		fmt.Println(fmt.Sprintf("\t%s: \n\t%s", a, fmt.Sprintf("----------------------------------------------------------------")))
		fmt.Println(fmt.Sprintf("\ttotal executions: \t%d", result.Total))
		fmt.Println(fmt.Sprintf("\tavg_duration: \t%vms\n", result.AverageTime))
	}
}

func printAllResponsesPerAssignment(r *scheduler.Result) {
	fmt.Println("\nall service responses per assignment: \n==================================")
	assignmentAndResponsesMap := r.Responses
	for a, r := range assignmentAndResponsesMap {
		fmt.Println(fmt.Sprintf("\t%s: \n\t----------------------------------------------------------------", a))
		statusCount := getAllResponsesPerAssignment(r)
		for s, t := range statusCount {
			fmt.Println(fmt.Sprintf("\tstatus: %s", s))
			fmt.Println(fmt.Sprintf("\ttotal: %d\n", t))
		}
	}
}

func getAllResponsesPerAssignment(responses []*apiclient.Response) map[string]int {
	statusCount := make(map[string]int)
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
