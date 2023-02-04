package util

import "os"

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
