package appdata

import (
	"encoding/json"
	"github.com/ghowland/sireus/code/util"
	"os"
	"time"
)

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
	Started       time.Time
	Finished      time.Time
	Score         float32
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
	Name string `json:"name"`
	Info string `json:"info"`

	VariableValues []BotVariableValue
	CommandHistory []ActionCommandResult
	LockTimers     []BotLockTimer
}

type BotQueryType int64

const (
	Range BotQueryType = iota
)

func (bqt BotQueryType) String() string {
	switch bqt {
	case Range:
		return "Range"
	}
	return "Unknown"
}

type BotQuery struct {
	QueryServer string       `json:"query_server"`
	QueryType   BotQueryType `json:"query_type"`
	Name        string       `json:"name"`
	Info        string       `json:"info"`
	Query       string       `json:"query"`
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
	Type     BotLockTimerType `json:"type"`
	Name     string           `json:"name"`
	Info     string           `json:"info"`
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
	Variables        []BotVariable             `json:"variables"`
	BotExtractor     BotExtractorQueryKey      `json:"bot_extractor"`
	States           []BotForwardSequenceState `json:"states"`
	LockTimers       []BotLockTimer            `json:"lock_timers"`
	BotTimeoutStale  time.Duration             `json:"bot_timeout_stale"`
	BotTimeoutRemove time.Duration             `json:"bot_timeout_remove"`
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
	Type                QueryServerType `json:"type"`
	Name                string          `json:"name"`
	Info                string          `json:"info"`
	Host                string          `json:"host"`
	Port                int             `json:"port"`
	AuthUser            string          `json:"auth_user"`
	AuthSecret          string          `json:"auth_secret"`
	DefaultDataDuration time.Duration   `json:"default_data_duration"`
}

type Site struct {
	Name          string        `json:"name"`
	Info          string        `json:"info"`
	BotGroupPaths []string      `json:"bot_group_paths"`
	QueryServers  []QueryServer `json:"query_servers"`
	BotGroups     []BotGroup
}

func LoadBotGroupConfig(path string) BotGroup {
	botGroupData, err := os.ReadFile(path)
	util.Check(err)

	var botGroup BotGroup
	err = json.Unmarshal(botGroupData, &botGroup)
	util.Check(err)

	return botGroup
}

func LoadSiteConfig(appConfig AppConfig) Site {
	siteData, err := os.ReadFile(appConfig.SiteConfigPath)
	util.Check(err)

	var site Site
	err = json.Unmarshal(siteData, &site)
	util.Check(err)

	// Load all our Bot Groups
	for _, botGroupPath := range site.BotGroupPaths {
		botGroup := LoadBotGroupConfig(botGroupPath)
		site.BotGroups = append(site.BotGroups, botGroup)
	}

	return site
}
