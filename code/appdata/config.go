package appdata

import (
	"encoding/json"
	"github.com/ghowland/sireus/code/util"
	"os"
)

type AppConfig struct {
	WebPath               string        `json:"web_path"`
	SiteConfigPath        string        `json:"site_config_path"`
	CurvePathFormat       string        `json:"curve_path_format"`
	QueryFastInternal     util.Duration `json:"query_fast_interval"` // BotQuery.Interval is overridden when users interact with the app, so they get fast interactive responses
	QueryFastDuration     util.Duration `json:"query_fast_duration"` // Duration QueryFastInterval is maintained after the last user interaction
	ReloadTemplatesAlways bool          `json:"reload_templates_always"`
	LogTemplateParsing    bool          `json:"log_template_parsing"`
}

func LoadConfig(path string) AppConfig {
	appConfigData, err := os.ReadFile(path)
	util.CheckPanic(err)

	var appConfig AppConfig
	err = json.Unmarshal(appConfigData, &appConfig)
	util.CheckPanic(err)

	return appConfig
}
