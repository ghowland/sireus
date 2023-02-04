package appdata

import (
	"encoding/json"
	"github.com/ghowland/sireus/code/util"
	"os"
)

type AppConfig struct {
	WebPath         string `json:"web_path"`
	ActionPath      string `json:"action_path"`
	CurvePathFormat string `json:"curve_path_format"`
}

func LoadConfig(path string) AppConfig {
	app_config_data, err := os.ReadFile(path)
	util.Check(err)

	var app_config AppConfig
	json.Unmarshal([]byte(app_config_data), &app_config)

	return app_config
}
