package data

import "time"

type (
	// Action is what is considered for execution.  It will receive a Final Score from it's Weight and Consideration Final Scores
	Action struct {
		Name               string                `json:"name"`
		Info               string                `json:"info"`
		IsLaunched         bool                  `json:"is_launched"`      // If false, this will never execute.  Launching means it is configured and ready to run live.  When Actions are created, is_launched==false and must be changed so that the action could execute.
		IsDisabled         bool                  `json:"is_disabled"`      // When testing changes, disable with modifying config
		Weight             float64               `json:"weight"`           // This is the multiplier for the Final Score, from the Consideration Final Score
		WeightMin          float64               `json:"weight_min"`       // If Weight != 0, then this is the Floor value.  We will bump it to this value, if it is less than this value
		WeightThreshold    float64               `json:"weight_threshold"` // If non-0, this is the threshold to be Active, and potentially execute Actions.  If the Final Score is less than this Threshold, this Action can never run.  WeightMin and WeightThreshold are independent tests, and will have different results when used together, so take that into consideration.
		RequiredLockTimers []string              `json:"required_lock_timers"`
		RequiredStates     []string              `json:"required_states"`
		Considerations     []ActionConsideration `json:"considerations"`
		Command            ActionCommand         `json:"command"`
	}
)

type (
	// Considerations are units for scoring an Action.  Each creates a Score, and they are combined to created the
	// Consideration Final Score.
	ActionConsideration struct {
		Name       string  `json:"name"`
		Weight     float64 `json:"weight"`
		CurveName  string  `json:"curve"`
		RangeStart float64 `json:"range_start"`
		RangeEnd   float64 `json:"range_end"`
		Evaluate   string  `json:"evaluate"`
	}
)

type (
	// What will we do with this ActionCommand?  We will only ever do 1 thing per Action, as this is a Decision System
	ActionCommandType int64
)

const (
	ShellCommand    ActionCommandType = iota // Shell Command
	WebHttps                                 // HTTPS request
	WebHttpInsecure                          // HTTP request
	WebRPC                                   // RPC call
	NoOperation                              // Do nothing
)

func (act ActionCommandType) String() string {
	switch act {
	case ShellCommand:
		return "ShellCommand"
	case WebHttps:
		return "WebHttps"
	case WebHttpInsecure:
		return "WebHttpInsecure"
	case WebRPC:
		return "WebRPC"
	case NoOperation:
		return "NoOperation"
	}
	return "Unknown"
}

type (
	// When an Action is selected for execution by it's Final Score, the ActionCommand is executed.  A command or web request.
	ActionCommand struct {
		Type              ActionCommandType `json:"type"`
		Content           string            `json:"content"`
		SuccessStatus     int               `json:"success_status"`
		SuccessContent    string            `json:"success_content"`
		LockTimerDuration Duration          `json:"lock_timer_duration"`
		HostExecKey       string            `json:"host_exec_key"`    // Sireus Client presents this key to get commands to run
		SetBotStates      []string          `json:"set_bot_states"`   // Will Advance all of these Bot States.  Advance can only go forward in the list, or start at the very beginning.  It can't go backwards, that is invalid data.
		JournalTemplate   string            `json:"journal_template"` // Templated Text formatted with variables from the Bot.VariableValues.  This is logged in JSON log-line and can be used to create Outage Reports, etc
	}
)

type (
	// When an Action is selected for execution by it's Final Score, the ActionCommand will execute and store this result.
	ActionCommandResult struct {
		ActionName    string
		ResultStatus  string
		ResultContent string
		HostExecOn    string
		Started       time.Time
		Finished      time.Time
		Score         float64
	}
)
