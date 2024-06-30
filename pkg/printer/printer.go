package printer

import (
	"fmt"

	"github.com/nadavbm/etzba/roles/apiclient"
	"github.com/nadavbm/etzba/roles/common"
)

// PrintToTerminal prints resutls to terminal
func PrintToTerminal(r *common.Result, collectApiResponses bool) {
	printAllTaskDurations(r)
	printDetailedAssignmentExecutions(r, collectApiResponses)
}

func printAllTaskDurations(r *common.Result) {
	fmt.Println("\nGeneral results: \n================")
	fmt.Println("")
	fmt.Println(fmt.Sprintf("\tjob duration: \t\t%v\t", r.General.JobDuration))
	fmt.Println(fmt.Sprintf("\trequest rate: \t\t%v/s\t", r.General.RequestRate))
	fmt.Println(fmt.Sprintf("\ttotal processed tasks:  %d\t", r.General.TotalExeuctions))
	fmt.Println(fmt.Sprintf("\tavg_duration: \t\t%.2fms", r.General.ProcessedDurations.AverageTime))
	fmt.Println(fmt.Sprintf("\tmin_duration: \t\t%.2fms", r.General.ProcessedDurations.MinimumTime))
	fmt.Println(fmt.Sprintf("\tmed_duration: \t\t%.2fms", r.General.ProcessedDurations.MedianTime))
	fmt.Println(fmt.Sprintf("\tmax_duration: \t\t%.2fms", r.General.ProcessedDurations.MaximumTime))
	return
}

func printDetailedAssignmentExecutions(r *common.Result, collectApiResponses bool) {
	fmt.Println("\nDetailed result per assignment: \n===============================")
	for _, a := range r.Assignments {
		fmt.Println(fmt.Sprintf("\n\t%s: \n\t%s", a.Title, fmt.Sprintf("-------------------------------------------------------------------------------------------------------------------------------------------------------------")))
		fmt.Println(fmt.Sprintf("\ttotal executions: \t%d", a.TotalExeuctions))
		if collectApiResponses {
			statusCount := getAllResponsesPerAssignment(a.ApiResponses)
			for s, t := range statusCount {
				fmt.Println(fmt.Sprintf("\tstatus: \t\t%s", s))
				fmt.Println(fmt.Sprintf("\ttotal requests: \t%d", t))
			}
		}
		fmt.Println(fmt.Sprintf("\tavg_duration: \t\t%.2fms", a.ProcessedDurations.AverageTime))
		fmt.Println(fmt.Sprintf("\tmin_duration: \t\t%.2fms", a.ProcessedDurations.MinimumTime))
		fmt.Println(fmt.Sprintf("\tmax_duration: \t\t%.2fms", a.ProcessedDurations.MaximumTime))
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
