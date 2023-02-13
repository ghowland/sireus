package data

type (
	// Web App server configuration
	AppConfig struct {
		WebPath                           string   `json:"web_path"`                             // Path to the Handlebars template content.  Holds *.hbs files
		SiteConfigPath                    string   `json:"site_config_path"`                     // Path to the config.yaml file that contains a Site.  For now only 1, but later will make this dynamic
		CurvePathFormat                   string   `json:"curve_path_format"`                    // String to format for each of the Curve JSON files, that contain the points we use to calculate from a curve
		ServerLoopDelay                   Duration `json:"server_loop_delay"`                    // After running the server loop, how long to delay, so we aren't in full spin lock.  This should be short like "0.8s"
		QueryLockTimeout                  Duration `json:"query_lock_timeout"`                   // We run Queries in the background, if they run longer than this, clear the lock.  This should be a longer time, like "60s".  TODO(ghowland): Pass in custom contexts and cancel them?  Better to really control it.
		QueryFastInternal                 Duration `json:"query_fast_interval"`                  // BotQuery.Interval is overridden when users interact with the app, so they get fast interactive responses
		QueryFastDuration                 Duration `json:"query_fast_duration"`                  // Duration QueryFastInterval is maintained after the last user interaction
		InteractiveSessionTimeout         Duration `json:"interactive_session_timeout"`          // Duration an InteractiveSession is kept until it is assumed finished, and can be purged
		InteractiveDurationMinutesDefault int      `json:"interactive_duration_minutes_default"` // How many minutes we default to starting the interactive query to.  15 minutes is reasonable
		PrometheusExportPort              int      `json:"prometheus_export_port"`               // Port used to listen for the Prometheus Exporter data we will put back into Prometheus.  For the main application, and the demo, if it is enabled
		EnableDemo                        bool     `json:"enable_demo"`                          // If true, the demo will be enabled and will export additional metrics to Prometheus to make learning Sireus easier.  The demo shares the PrometheusExportPort for simplicity
		ReloadTemplatesAlways             bool     `json:"reload_templates_always"`              // For development, if true this will always reload Handlebars files.  Only need to restart the server to rebuild code.
		LogTemplateParsing                bool     `json:"log_template_parsing"`                 // For web development debugging, if true this will print out all the templates that are parsed.  It's not generally useful, but if you are having a problem with Handlebars template imports or related it can help
	}
)
