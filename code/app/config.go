package app

import (
	"encoding/json"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/util"
	"os"
)

// Load the Server config
func LoadConfig(path string) data.AppConfig {
	appConfigData, err := os.ReadFile(path)
	util.CheckPanic(err)

	var appConfig data.AppConfig
	err = json.Unmarshal(appConfigData, &appConfig)
	util.CheckPanic(err)

	return appConfig
}
