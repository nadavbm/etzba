package printer

import (
	"fmt"

	"github.com/nadavbm/etzba/pkg/calculator"
	"github.com/nadavbm/etzba/roles/apiclient"
	"github.com/nadavbm/etzba/roles/scheduler"
)

// PrintToTerminal prints resutls to terminal
func PrintToTerminal(r *scheduler.Result, collectApiResponses bool) {
	printAllTaskDurations(r)
	printDetailedAssignmentExecutions(r, collectApiResponses)
}

func printAllTaskDurations(r *scheduler.Result) {
	fmt.Println("\nGeneral results: \n================")
	fmt.Println("")
	if r.JobDuration != 0 {
		fmt.Println(fmt.Sprintf("\tjob duration: \t\t%v\t", r.JobDuration))
	}
	fmt.Println(fmt.Sprintf("\trequest rate: \t\t%d/s\t", r.RequestRate))
	c := calculator.NewCalculator()
	result := c.GetResult(r.Durations)
	fmt.Println(fmt.Sprintf("\ttotal processed tasks:  %d\t", result.Total))
	fmt.Println(fmt.Sprintf("\ttotal processing time:\t%vms", result.TotalJobTime))
	fmt.Println(fmt.Sprintf("\tavg_duration: \t\t%vms", result.AverageTime))
	fmt.Println(fmt.Sprintf("\tmin_duration: \t\t%vms", result.MinimumTime))
	fmt.Println(fmt.Sprintf("\tmed_duration: \t\t%vms", result.MedianTime))
	fmt.Println(fmt.Sprintf("\tmax_duration: \t\t%vms", result.MaximumTime))
	return
}

func printDetailedAssignmentExecutions(r *scheduler.Result, collectApiResponses bool) {
	fmt.Println("\nDetailed result per assignment: \n===============================")

	for a, d := range r.Assignments {
		c := calculator.NewCalculator()
		result := c.GetResult(d)
		fmt.Println(fmt.Sprintf("\n\t%s: \n\t%s", a, fmt.Sprintf("------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")))
		fmt.Println(fmt.Sprintf("\ttotal executions: \t%d", result.Total))
		if collectApiResponses {
			statusCount := getAllResponsesPerAssignment(r.Responses[a])
			for s, t := range statusCount {
				fmt.Println(fmt.Sprintf("\tstatus: \t\t%s", s))
				fmt.Println(fmt.Sprintf("\ttotal requests: \t%d", t))
			}
		}
		fmt.Println(fmt.Sprintf("\tavg_duration: \t\t%vms", result.AverageTime))
		fmt.Println(fmt.Sprintf("\tmin_duration: \t\t%vms", result.MinimumTime))
		fmt.Println(fmt.Sprintf("\tmax_duration: \t\t%vms", result.MaximumTime))
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
