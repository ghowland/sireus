package appdata

import (
	"encoding/json"
	"github.com/ghowland/sireus/code/util"
	"os"
)

type AppConfig struct {
	WebPath         string `json:"web_path"`
	SiteConfigPath  string `json:"site_config_path"`
	CurvePathFormat string `json:"curve_path_format"`
}

func LoadConfig(path string) AppConfig {
	appConfigData, err := os.ReadFile(path)
	util.Check(err)

	var appConfig AppConfig
	err = json.Unmarshal(appConfigData, &appConfig)
	util.Check(err)

	return appConfig
}
