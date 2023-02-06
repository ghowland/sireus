package appdata

import (
	"github.com/ghowland/sireus/code/util"
	"time"
)

type ActionConsideration struct {
	Name       string  `json:"name"`
	Weight     float32 `json:"weight"`
	CurveName  string  `json:"curve"`
	RangeStart float32 `json:"range_start"`
	RangeEnd   float32 `json:"range_end"`
	Evaluate   string  `json:"evaluate"`
}

type ActionCommandType int64

const (
	Bash ActionCommandType = iota
	WebHttps
	WebHttpInsecure
	WebRPC
)

func (act ActionCommandType) String() string {
	switch act {
	case Bash:
		return "Bash"
	case WebHttps:
		return "WebHttps"
	case WebHttpInsecure:
		return "WebHttpInsecure"
	case WebRPC:
		return "WebRPC"
	}
	return "Unknown"
}

type ActionCommandResult struct {
	ActionName    string
	ResultStatus  string
	ResultContent string
	HostExecOn    string
	Started       time.Time
	Finished      time.Time
	Score         float32
}

type ActionCommand struct {
	Type              ActionCommandType `json:"type"`
	Content           string            `json:"content"`
	SuccessStatus     int               `json:"success_status"`
	SuccessContent    string            `json:"success_content"`
	LockTimerDuration util.Duration     `json:"lock_timer_duration"`
	HostExecKey       string            `json:"host_exec_key"` // Sireus Client presents this key to get commands to run
}

type Action struct {
	Name               string                `json:"name"`
	Info               string                `json:"info"`
	IsDisabled         bool                  `json:"is_disabled"`      // When testing changes, disable with modifying config
	Weight             float32               `json:"weight"`           // This is the multiplier for the Final Score, from the Consideration Final Score
	WeightMin          float32               `json:"weight_min"`       // If Weight != 0, then this is the Floor value.  We will bump it to this value, if it is less than this value
	WeightThreshold    float32               `json:"weight_threshold"` // If non-0, this is the threshold to be Active, and potentially execute Actions.  If the Final Score is less than this Threshold, this Action can never run.  WeightMin and WeightThreshold are independent tests, and will have different results when used together, so take that into consideration.
	RequiredLockTimers []string              `json:"required_lock_timers"`
	RequiredStates     []string              `json:"required_states"`
	SetBotStates       []string              `json:"set_bot_states"`
	Considerations     []ActionConsideration `json:"considerations"`
	Command            ActionCommand         `json:"command"`
}

type BotActionData struct {
	ActionName             string          // Action.Name matches to store data about that action per Bot.  Can use a map[string]BotActionData
	FinalScore             bool            // Final Score is the total result of calculations to Score this action for execution
	IsActive               bool            // This Action is Active is the FinalScore is over the WeightThreshold, even if it is not executed
	ActiveStartTime        time.Time       // Time this Active started, so we can use it for an Evaluation variable
	LastExecutedActionTime time.Time       // Last time we executed this Action
	Time                   time.Time       // When this was updated
	History                []BotActionData // We keep N records history, but no recursive depth.  Top level keeps history, no history nodes keep history
}

type BotVariableValue struct {
	Name  string
	Value float32
	Time  time.Time
}

type Bot struct {
	Name           string
	VariableValues []BotVariableValue
	StateValues    []string
	CommandHistory []ActionCommandResult
	LockTimers     []BotLockTimer
	ActionData     map[string]BotActionData // Key is Action.Name
}

type BotQueryType int64

const (
	Range BotQueryType = iota
)

func (bqt BotQueryType) String() string {
	switch bqt {
	case Range:
		return "query_range"
	}
	return "Unknown"
}

type BotQuery struct {
	QueryServer string        `json:"query_server"`
	QueryType   BotQueryType  `json:"query_type"`
	Name        string        `json:"name"`
	Info        string        `json:"info"`
	Query       string        `json:"query"`
	Interval    util.Duration `json:"interval"`
}

type BotForwardSequenceState struct {
	Name   string   `json:"name"`
	Info   string   `json:"info"`
	Labels []string `json:"labels"`
}

type BotExtractorQueryKey struct {
	QueryName string `json:"query_name"`
	Key       string `json:"key"`
}

type BotLockTimerType int64

const (
	LockBotGroup BotLockTimerType = iota
	LockBot
)

func (bltt BotLockTimerType) String() string {
	switch bltt {
	case LockBotGroup:
		return "Lock Bot Group"
	case LockBot:
		return "Lock Bot"
	}
	return "Unknown"
}

type BotLockTimer struct {
	Type           BotLockTimerType `json:"type"`
	Name           string           `json:"name"`
	Info           string           `json:"info"`
	IsActive       bool
	Timeout        time.Time
	ActivatedByBot string // Bot.Name of who set this Lock Timer, so we can track Actions
}

type BotVariableType int64

const (
	Boolean BotVariableType = iota
	Float
)

func (bvt BotVariableType) String() string {
	switch bvt {
	case Boolean:
		return "Boolean"
	case Float:
		return "Float"
	}
	return "Unknown"
}

type BotVariable struct {
	Type           BotVariableType `json:"type"`
	Name           string          `json:"name"`
	BotKey         string          `json:"bot_key"` // Determines which Metric matches a Bot, may change between queries
	QueryName      string          `json:"query_name"`
	QueryKey       string          `json:"query_key"`       // Metric key to extract
	QueryKeyValue  string          `json:"query_key_value"` // Metric key value to match against the QueryKey
	Evaluate       string          `json:"evaluate"`        // If this is non-empty, query will not be performed.  After query testing for other variables, this will have a final phase of processing, and will take all the query-made variables and perform govaluation.Evaluate() with this evaluate string, to set this variable.  Evaluate variables cannot use each other, only Query variables.
	BoolRangeStart float32         `json:"bool_range_start"`
	BoolRangeEnd   float32         `json:"bool_range_end"`
	BoolInvert     bool            `json:"bool_invert"`
	Export         bool            `json:"export"` // If true, this variable will be exported for Metric collection.  Normally not useful, because we just got it from the Metric system.
}

type BotGroup struct {
	Name             string                    `json:"name"`
	Info             string                    `json:"info"`
	Queries          []BotQuery                `json:"queries"`
	Variables        []BotVariable             `json:"variables"`
	BotExtractor     BotExtractorQueryKey      `json:"bot_extractor"`
	States           []BotForwardSequenceState `json:"states"`
	LockTimers       []BotLockTimer            `json:"lock_timers"`
	BotTimeoutStale  util.Duration             `json:"bot_timeout_stale"`
	BotTimeoutRemove util.Duration             `json:"bot_timeout_remove"`
	ActionScoreMin   float32                   `json:"action_score_min"` // Minimum score to execute Action
	Actions          []Action                  `json:"actions"`
	Bots             []Bot

	// Invalid = Isn't getting all the information.  Stale = Information out of data.  Removed = No data for too long, removing.
	InvalidBots []string
	StaleBots   []string
	RemovedBots []string
}

type QueryServerType int64

const (
	Prometheus QueryServerType = iota
)

func (qst QueryServerType) String() string {
	switch qst {
	case Prometheus:
		return "Prometheus"
	}
	return "Unknown"
}

type QueryServer struct {
	ServerType          QueryServerType `json:"server_type"`
	Name                string          `json:"name"`
	Info                string          `json:"info"`
	Host                string          `json:"host"`
	Port                int             `json:"port"`
	AuthUser            string          `json:"auth_user"`
	AuthSecret          string          `json:"auth_secret"`
	DefaultStep         string          `json:"default_step"`
	DefaultDataDuration util.Duration   `json:"default_data_duration"`
	WebUrlFormat        string          `json:"web_url_format"`
}

type Site struct {
	Name          string        `json:"name"`
	Info          string        `json:"info"`
	BotGroupPaths []string      `json:"bot_group_paths"`
	QueryServers  []QueryServer `json:"query_servers"`
	BotGroups     []BotGroup
}
