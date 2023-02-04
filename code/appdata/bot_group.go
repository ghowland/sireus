package appdata

import "time"

type ActionConsideration struct {
	Name       string  `json:"name"`
	Weight     float32 `json:"weight"`
	CurveName  string  `json:"curve"`
	RangeStart float32 `json:"range_start"`
	RangeEnd   float32 `json:"range_end"`
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
}

type ActionCommand struct {
	Type           ActionCommandType `json:"type"`
	Content        string            `json:"content"`
	SuccessStatus  int               `json:"success_status"`
	SuccessContent string            `json:"success_content"`
	HostExecKey    string            `json:"host_exec_key"`
}

type Action struct {
	Name           string                `json:"name"`
	Info           string                `json:"info"`
	Weight         float32               `json:"weight"`
	WeightMin      float32               `json:"weight_min"`
	Command        ActionCommand         `json:"command"`
	Considerations []ActionConsideration `json:"considerations"`
}

type BotVariableValue struct {
	Name  string
	Value float32
	Time  time.Time
}

type Bot struct {
	Name    string   `json:"name"`
	Info    string   `json:"info"`
	Actions []Action `json:"actions"`

	VariableValues []BotVariableValue
}

type BotQueryType int64

const (
	PrometheusQueryRange BotQueryType = iota
)

func (bqt BotQueryType) String() string {
	switch bqt {
	case PrometheusQueryRange:
		return "PromQueryRange"
	}
	return "Unknown"
}

type BotQuery struct {
	QueryServer string `json:"query_server"`
	Name        string `json:"name"`
	Info        string `json:"info"`
	Query       string `json:"query"`
}

type BotForwardSequenceState struct {
	Name   string   `json:"name"`
	Info   string   `json:"info"`
	Labels []string `json:"labels"`
}

type BotExtractorQueryKey struct {
	QueryName string `json:"query_name"`
	Keys      string `json:"key"`
}

type BotLockTimer struct {
	Name     string `json:"name"`
	IsActive bool
	Timeout  time.Time
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
	QueryName      string          `json:"query_name"`
	QueryKey       string          `json:"query_key"`
	BoolRangeStart float32         `json:"bool_range_start"`
	BoolRangeEnd   float32         `json:"bool_range_end"`
	BoolInvert     bool            `json:"bool_invert"`
}

type BotGroup struct {
	Name             string                    `json:"name"`
	Info             string                    `json:"info"`
	Queries          []BotQuery                `json:"queries"`
	BotExtractor     BotExtractorQueryKey      `json:"bot_extractor"`
	States           []BotForwardSequenceState `json:"states"`
	Actions          []Action                  `json:"actions"`
	Variables        []BotVariable             `json:"variables"`
	LockTimers       []BotLockTimer            `json:"lock_timers"`
	BotTimeoutStale  time.Duration             `json:"bot_timeout_stale"`
	BotTimeoutRemove time.Duration             `json:"bot_timeout_remove"`
	Bots             []Bot

	// Invalid = Isn't getting all the information.  Stale = Information out of data.  Removed = No data for too long, removing.
	InvalidBots []string
	StaleBots   []string
	RemovedBots []string
}

type BotQueryServer struct {
	Type         BotQueryType  `json:"type"`
	Name         string        `json:"name"`
	Info         string        `json:"info"`
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	Auth         string        `json:"auth"`
	DataDuration time.Duration `json:"data_duration"`
}

type Site struct {
	Name         string           `json:"name"`
	Paths        []string         `json:"paths"`
	QueryServers []BotQueryServer `json:"query_servers"`
	BotGroups    []BotGroup
}
