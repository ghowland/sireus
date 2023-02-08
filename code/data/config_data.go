package data

type AppConfig struct {
	WebPath               string   `json:"web_path"`
	SiteConfigPath        string   `json:"site_config_path"`
	CurvePathFormat       string   `json:"curve_path_format"`
	QueryFastInternal     Duration `json:"query_fast_interval"` // BotQuery.Interval is overridden when users interact with the app, so they get fast interactive responses
	QueryFastDuration     Duration `json:"query_fast_duration"` // Duration QueryFastInterval is maintained after the last user interaction
	ReloadTemplatesAlways bool     `json:"reload_templates_always"`
	LogTemplateParsing    bool     `json:"log_template_parsing"`
}
