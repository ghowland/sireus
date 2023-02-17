package util

import (
	"fmt"
	"github.com/ghowland/sireus/code/data"
	"time"
)

// Returns the Interaction Start Time in the Format required by the HTML datetime picker
func FormatInteractiveStartTime() string {
	// 15 minutes ago
	//TODO(ghowland): Remove hard-code, put into AppConfig, also make default Duration in the webapp
	var t = GetTimeNow().Add(time.Duration(-data.SireusData.AppConfig.InteractiveDurationMinutesDefault*60) * time.Second)

	ampm := "AM"
	hour := t.Hour()
	if hour > 12 {
		hour -= 12
		ampm = "PM"
	}

	output := fmt.Sprintf("%02d/%02d/%d, %d:%02d %s", t.Day(), t.Month(), t.Year(), hour, t.Minute(), ampm)
	return output
}
