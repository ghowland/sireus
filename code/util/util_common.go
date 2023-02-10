package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aymerick/raymond"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

// Call Check when we only want to log the error and wrap error testing, but it does not require an exceptional response
func Check(e error) bool {
	if e != nil {
		log.Printf("ERROR: %s", e.Error())
		return true
	}

	return false
}

// Call CheckNoLog when we want to get a boolean on the error, but dont want to log because we handle the response and
// it's too noisy or not useful.
func CheckNoLog(e error) bool {
	if e != nil {
		return true
	}

	return false
}

// Call CheckPanic for configuration errors that can't be solved.
func CheckPanic(e error) {
	if e != nil {
		log.Printf("PANIC: %s", e.Error())
		panic(e)
	}
}

// Does this file exist?  Wrap to make code shorter in situ
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Format Handlebars string, so I don't have to remember any arguments
func HandlebarFormatText(format string, mapData map[string]string) string {
	result, err := raymond.Render(format, mapData)
	Check(err)

	return result
}

// Test if a string is in a slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

var (
	// Float64 is our primary data type, custom error for tracking problems
	InvalidTypeFloat64 = errors.New("Value could not be converted to Float64")
)

// Convert any value we can into a float64 in a predictable manner
func ConvertInterfaceToFloat(value interface{}) (float64, error) {
	switch value := value.(type) {
	case float64:
		return value, nil
	case float32:
		return float64(value), nil
	case int64:
		return float64(value), nil
	case int32:
		return float64(value), nil
	case int:
		return float64(value), nil
	case uint64:
		return float64(value), nil
	case uint32:
		return float64(value), nil
	case uint:
		return float64(value), nil
	default:
		return math.NaN(), InvalidTypeFloat64
	}
}

// Clamp a value between a min and a max
func Clamp(value float64, min float64, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}

// Returns clamped value between 0-1, where the value falls between the range
func RangeMapper(value float64, rangeMin float64, rangeMax float64) float64 {
	valueRange := rangeMax - rangeMin

	rawRangeValue := value / valueRange

	finalValue := Clamp(rawRangeValue, 0, 1)

	return finalValue
}

// Converts a boolean to a string of "0" or "1"
func BoolToFloatString(value bool) string {
	if value {
		return "1"
	} else {
		return "0"
	}
}

// Print JSON, for debugging
func PrintJson(value interface{}) string {
	output, err := json.MarshalIndent(value, "", "  ")
	Check(err)

	return string(output)
}

// Print a string array.  For human readability or debugging
func PrintStringArrayCSV(slice []string) string {
	output := strings.Join(slice, ", ")

	return string(output)
}

func FormatTimeLong(t time.Time) string {
	utc := t.UTC()

	output := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", utc.Year(), utc.Month(), utc.Day(), utc.Hour(), utc.Minute(), utc.Second())

	return output
}
