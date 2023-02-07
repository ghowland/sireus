package appdata

import (
	"fmt"
	"github.com/BenJetson/humantime"
	"time"
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
		return humanize.Ordinal(int64(value))
	case FormatComma:
		return humanize.Comma(int64(value))
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
