package appdata

import (
	"fmt"
	"github.com/BenJetson/humantime"
	"time"
	"github.com/dustin/go-humanize"
)

func FormatBotVariable(format BotVariableFormat, value float64) string {
	switch format {
	case FormatFloat:
		return fmt.Sprintf("%.2f", value)
	case FormatBytes:
		return humanize.Bytes(uint64(value))
	case FormatBandwidth:
		return humanize.Bytes(uint64(value))
	case FormatDuration:
		return humantime.Duration(time.Duration(value))
	case FormatTime:
		return humanize.Time(time.Unix(int64(value), 0))
	case FormatOrdinal:
		return humanize.Ordinal(int(value))
	case FormatComma:
		return humanize.Comma(int64(value))
	case FormatMetricPrefix:
		return humanize.SI(value, "")
	case FormatPercent:
		return fmt.Sprintf("%.1f%%", value*100)
	case FormatBool:
		if value == 0 {
			return "False"
		} else {
			return "True"
		}
	default:
		return fmt.Sprintf("%.2f", value)
	}
}
