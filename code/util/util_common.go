package util

import (
	"github.com/aymerick/raymond"
	"os"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func HandlebarFormatText(format string, mapData map[string]string) string {
	result, err := raymond.Render(format, mapData)
	Check(err)

	return result
}
