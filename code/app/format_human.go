package app

import (
	"fmt"
	"github.com/BenJetson/humantime"
	"github.com/dustin/go-humanize"
	"github.com/ghowland/sireus/code/data"
	"time"
)

// Bot.VariableValues are all floats, but we want them to have human readable strings
func FormatBotVariable(format data.BotVariableFormat, value float64) string {
	switch format {
	case data.FormatFloat:
		return fmt.Sprintf("%.2f", value)
	case data.FormatBytes:
		return humanize.Bytes(uint64(value))
	case data.FormatBandwidth:
		return humanize.Bytes(uint64(value))
	case data.FormatDuration:
		return humantime.Duration(time.Duration(value))
	case data.FormatTime:
		return humanize.Time(time.Unix(int64(value), 0))
	case data.FormatOrdinal:
		return humanize.Ordinal(int(value))
	case data.FormatComma:
		return humanize.Comma(int64(value))
	case data.FormatMetricPrefix:
		return humanize.SI(value, "")
	case data.FormatPercent:
		return fmt.Sprintf("%.1f%%", value*100)
	case data.FormatBool:
		if value == 0 {
			return "False"
		} else {
			return "True"
		}
	case data.FormatRequestsPerSecond:
		return fmt.Sprintf("%.0f/s", value)
	case data.FormatInteger:
		return fmt.Sprintf("%.0f", value)
	default:
		return fmt.Sprintf("%.2f", value)
	}
}
