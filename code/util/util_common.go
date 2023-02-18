package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aymerick/raymond"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Call Check when we want to get a boolean on the error, but dont want to log because we handle the response and it's too noisy or not useful.
func Check(e error) bool {
	if e != nil {
		return true
	}

	return false
}

// Call CheckLog when we only want to log the error and wrap error testing, but it does not require an exceptional response
func CheckLog(e error) bool {
	if e != nil {
		log.Printf("ERROR: %s", e.Error())
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

func FileLoad(path string) (string, error) {
	dat, err := os.ReadFile(path)
	if Check(err) {
		return "", err
	}

	return string(dat), nil
}

// Format Handlebars string from strings, so I don't have to remember any arguments
func HandlebarFormatText(format string, mapData map[string]string) string {
	result, err := raymond.Render(format, mapData)
	Check(err)

	return result
}

func HandlebarsRegisterPartials(partialsDirPrefix string, removeBasePath string, template *raymond.Template) {
	topPaths, err := filepath.Glob(partialsDirPrefix)
	Check(err)

	// Import partials into handlebars template
	for _, topPath := range topPaths {
		fileInfo, err := os.Stat(topPath)
		if Check(err) {
			log.Printf(fmt.Sprintf("Could not stat path: %s  Error: %s", topPaths, err))
			continue
		}

		if strings.HasSuffix(topPath, ".hbs") {
			partialName := strings.Replace(topPath, ".hbs", "", 1)
			partialName = strings.Replace(partialName, removeBasePath, "", 1)
			content, _ := FileLoad(topPath)
			template.RegisterPartial(partialName, content)
		} else if fileInfo.IsDir() {
			nextPath := fmt.Sprintf("%s/*", topPath)
			HandlebarsRegisterPartials(nextPath, removeBasePath, template)
		}
	}
}

// Format Handlebars string from data, so I don't have to remember any arguments
func HandlebarFormatData(format string, mapData map[string]interface{}) string {
	template := raymond.MustParse(format)

	HandlebarsRegisterPartials("web/partials/*", "web/", template)

	result := template.MustExec(mapData)

	return result
}

// Test if a string is in a slice
func StringInSlice(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func StringSliceRemoveString(slice []string, remove string) ([]string, error) {
	removeIndex, err := StringSliceFindIndex(slice, remove)
	if Check(err) {
		return slice, err
	}

	newSlice := StringSliceRemoveIndex(slice, removeIndex)
	return newSlice, nil
}

// Remove a string at index.  Wrapper for go-ness
func StringSliceRemoveIndex(slice []string, index int) []string {
	slice = append(slice[:index], slice[index+1:]...)
	return slice
}

// Returns the index of a string in a slice
func StringSliceFindIndex(slice []string, find string) (int, error) {
	for index, value := range slice {
		if value == find {
			return index, nil
		}
	}

	return -1, errors.New(fmt.Sprintf("Couldn't find string in slice: %s", find))
}

// Clamp a value between a min and a max
func Clamp(value float64, min float64, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}

// Returns clamped value between 0-1, where the value falls between the range
func RangeMapper(value float64, rangeMin float64, rangeMax float64) float64 {
	var valueRange float64
	var rawValue float64

	// Are we doing this backwards, going from high to low?
	isDecreaing := false

	if rangeMax == rangeMin {
		// Avoid division by zero, and just say it is always true
		return 1
	} else if rangeMax > rangeMin {
		// Get the distance between the values
		valueRange = rangeMax - rangeMin
	} else {
		isDecreaing = true

		// Get the distance between the values
		valueRange = rangeMin - rangeMax
	}

	// Proportional Value Range
	rawValue = value / valueRange

	// Clamp between 0 and 1
	finalValue := Clamp(rawValue, 0, 1)

	// If we are going the opposite direction, invert the value
	if isDecreaing {
		finalValue = 1 - finalValue
	}

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

// Print JSON, for transport
func PrintJsonData(value interface{}) string {
	output, err := json.Marshal(value)
	Check(err)

	return string(output)
}

// Print a string array.  For human readability or debugging
func PrintStringArrayCSV(slice []string) string {
	output := strings.Join(slice, ", ")

	return string(output)
}

// Format the time in ISO 8601, without the millis
func FormatTimeLong(t time.Time) string {
	utc := t.UTC()

	output := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", utc.Year(), utc.Month(), utc.Day(), utc.Hour(), utc.Minute(), utc.Second())

	return output
}

// Returns time.Now() in UTC.  Convenience wrapper, so it's never forgotten, because everything must always be in UTC
func GetTimeNow() time.Time {
	return time.Now().UTC()
}

// Replace any characters in unsafeChars with the replace string.  Quickly convert into a safe string
func StringReplaceUnsafeChars(value string, unsafeChars string, replace string) string {
	for _, unsafeChar := range unsafeChars {
		value = strings.Replace(value, string(unsafeChar), replace, -1)
	}
	return value
}

// Copy a string slice, because direct assignment is a reference
func CopyStringSlice(src []string) []string {
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}

func HttpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	bodyStr := string(body)
	return bodyStr, nil
}
