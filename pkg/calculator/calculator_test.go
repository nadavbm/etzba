package calculator

import (
	"testing"
	"time"
)

func TestCalculations(t *testing.T) {
	allDurations := []time.Duration{56625536, 103485052, 145246835, 143142433, 158100141, 150688596, 150993677, 165231877, 168389116, 158563818, 162723248, 179804452, 175139734, 193645495, 192880881, 188678533, 198230416, 204657193, 195815059, 198757912, 210586182, 213098556, 193328058, 201342432}

	min := GetMinimumTime(allDurations)
	if min != 56.625536 {
		t.Error("expected minimum duration to be 56.625536ms instead got", min)
	}

	avg := GetAverageTime(allDurations)
	if avg != 171.21480133333338 {
		t.Error("expected average duration to be 171.21480133333338ms instead got", avg)
	}

	med := GetMedianTime(allDurations)
	if med != 175.139734 {
		t.Error("expected median duration to be 175.139734ms instead got", med)
	}

	max := GetMaximumTime(allDurations)
	if max != 213.098556 {
		t.Error("expected maximum duration to be 213.098556ms instead got", max)
	}
}
