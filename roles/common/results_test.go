package common

import (
	"testing"
	"time"
)

func TestAppendAssignmentDurationsToConcatDurations(t *testing.T) {
	getDurations := []time.Duration{
		time.Duration(13.123 * float64(time.Millisecond)),
		time.Duration(12.766 * float64(time.Millisecond)),
		time.Duration(14.234 * float64(time.Millisecond)),
	}
	postDurations := []time.Duration{
		time.Duration(12.123 * float64(time.Millisecond)),
		time.Duration(13.523 * float64(time.Millisecond)),
		time.Duration(13.765 * float64(time.Millisecond)),
	}

	titles := []string{"URL: http://localhost:8080/v1/results, Method: GET", "URL: http://localhost:8080/v1/results, Method: POST"}

	allAssignmentsExecutions := make(map[string][]time.Duration)
	allAssignmentsExecutions[titles[0]] = getDurations
	allAssignmentsExecutions[titles[1]] = postDurations

	if allAssignmentsExecutions[titles[0]][0] != time.Duration(13.123*float64(time.Millisecond)) {
		t.Error("expected 13.123ms got ", allAssignmentsExecutions[titles[0]][0])
	}

	if allAssignmentsExecutions[titles[0]][1] != time.Duration(12.766*float64(time.Millisecond)) {
		t.Error("expected 12.766ms got ", allAssignmentsExecutions[titles[0]][0])
	}

	if allAssignmentsExecutions[titles[0]][2] != time.Duration(14.234*float64(time.Millisecond)) {
		t.Error("expected 14.234ms got ", allAssignmentsExecutions[titles[0]][0])
	}

	if allAssignmentsExecutions[titles[1]][0] != time.Duration(12.123*float64(time.Millisecond)) {
		t.Error("expected 12.123ms got ", allAssignmentsExecutions[titles[0]][0])
	}

	allDurations := concatAllDurations(allAssignmentsExecutions)
	if allDurations[0] != time.Duration(13.123*float64(time.Millisecond)) {
		t.Error("expected 13.123ms got ", allDurations[0])
	}

	if allDurations[1] != time.Duration(12.766*float64(time.Millisecond)) {
		t.Error("expected 12.766ms got ", allDurations[1])
	}

	if allDurations[2] != time.Duration(14.234*float64(time.Millisecond)) {
		t.Error("expected 14.234ms got ", allDurations[2])
	}

	if allDurations[3] != time.Duration(12.123*float64(time.Millisecond)) {
		t.Error("expected 12.123ms got ", allDurations[3])
	}

	if allDurations[4] != time.Duration(13.523*float64(time.Millisecond)) {
		t.Error("expected 13.523ms got ", allDurations[4])
	}
}
