package appdata

import (
	"github.com/ghowland/sireus/code/util"
	"time"
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

type ActionCommandType int64

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

type (
	// When an Action is selected for execution by it's Final Score, the ActionCommand is executed.  A command or web request.
	ActionCommand struct {
		Type              ActionCommandType `json:"type"`
		Content           string            `json:"content"`
		SuccessStatus     int               `json:"success_status"`
		SuccessContent    string            `json:"success_content"`
		LockTimerDuration util.Duration     `json:"lock_timer_duration"`
		HostExecKey       string            `json:"host_exec_key"`    // Sireus Client presents this key to get commands to run
		SetBotStates      []string          `json:"set_bot_states"`   // Will Advance all of these Bot States.  Advance can only go forward in the list, or start at the very beginning.  It can't go backwards, that is invalid data.
		JournalTemplate   string            `json:"journal_template"` // Templated Text formatted with variables from the Bot.VariableValues.  This is logged in JSON log-line and can be used to create Outage Reports, etc
	}
)

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
	// This stores the Final Scores and related data for all Actions, so they can be compared to determin if any
	// Action should be executed
	BotActionData struct {
		FinalScore                   float64            // Final Score is the total result of calculations to Score this action for execution
		IsAvailable                  bool               // This Action is Available (not blocked) if the FinalScore is over the WeightThreshold
		AvailableStartTime           time.Time          // Time IsAvailable started, so we can use it for an internal Evaluation variable "_available_start_time".  Stateful.
		LastExecutedActionTime       time.Time          // Last time we executed this Action.  Stateful.
		Details                      []string           // Details about the Evaluation and Scoring, to make it easier to understand the result
		ConsiderationFinalScores     map[string]float64 // Considerations Final Results for this Bot
		ConsiderationEvaluatedScores map[string]float64 // Considerations Evaluated score, but not weighted results for this Bot
	}
)

type (
	// Bots the core structure for this system.  They are ephemeral and build from the Bot Group data, and store
	// minimal data.  Bots are expected to be added or removed at any time, and there is a Timeout for Stale, Invalid,
	// and Removed bots.
	//
	// All Bots are expected to get all the data specified from the Bot Group in their Query to
	// Variable mapping.
	//
	// If a Bot is missing any data for it's variables, it is considered Invalid, because we are not
	// operating with a full set of data.
	Bot struct {
		Name                 string
		VariableValues       map[string]float64
		SortedVariableValues util.PairFloat64List // Sorted VariableValues, Handlebars helper
		StateValues          []string
		CommandHistory       []ActionCommandResult
		LockTimers           []BotLockTimer
		ActionData           map[string]BotActionData // Key is Action.Name
		SortedActionData     PairBotActionDataList    // Scored ActionData, Handlebars helper
		FreezeActions        bool                     // If true, no actions will be taken for this Bot.  Single agent control
	}
)

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

type (
	// These queries are stored in BotGroup, but are used to populate the Bots with query Variables.
	BotQuery struct {
		QueryServer string        `json:"query_server"`
		QueryType   BotQueryType  `json:"query_type"`
		Name        string        `json:"name"`
		Info        string        `json:"info"`
		Query       string        `json:"query"`
		Interval    util.Duration `json:"interval"`
	}
)

type (
	// Forward Sequence State is the term I am using to describe a State Machine that only has a single forward
	// sequence.  It can be Advanced and it can be Reset, but the state cannot go backwards.
	//
	// In this way you can create a State Machine for investigating problems, trying to solve them, checking for
	// resolution, and finally escalating and waiting for someone else to fix it.
	//
	// If a resolution is detected by an Action, the action can Reset this state, starting the States process over again.
	//
	// States are used to exclude Actions from being tested, so that Actions can be targetted at a specific State of a
	// Bot's operation.  This allows segmenting logic.  Actions use Action.RequiredStates to limit when they can execute.
	BotForwardSequenceState struct {
		Name   string   `json:"name"`
		Info   string   `json:"info"`
		Labels []string `json:"labels"`
	}
)

type (
	// This is how Bots are created.  There is a BotQuery named QueryName that will use the Key to find the name of the
	// Bots.  Using something like "instance", "node" or "service" is recommended, that will uniquely identify a Bot
	// inside a BotGroup's configuration.
	BotExtractorQueryKey struct {
		QueryName string `json:"query_name"`
		Key       string `json:"key"`
	}
)

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

type (
	// BotLockTimer is used to both block an Action from being executed, if the BotLockTimer.IsActive and has not
	// reached the Timeout yet.  Actions can use multiple BotLockTimers which essentially act as execution "channels"
	// where Actions execute 1 at a time.
	//
	// BotLockTimeType specifies the scope of the lock.  Is it locked at the Bot level or the BotGroup level?
	// BotGroup level locks (LockBotGroup) are essentially global level locks, as BotGroups do not interact with each
	// other, as they are data silos for decision-making.
	BotLockTimer struct {
		Type           BotLockTimerType `json:"type"`
		Name           string           `json:"name"`
		Info           string           `json:"info"`
		IsActive       bool
		Timeout        time.Time
		ActivatedByBot string // Bot.Name of who set this Lock Timer, so we can track Actions
	}
)

type BotVariableType int64

const (
	Boolean BotVariableType = iota
	Float
)

type BotVariableFormat int64

const (
	FormatFloat BotVariableFormat = iota
	FormatBool
	FormatBytes
	FormatBandwidth
	FormatTime
	FormatDuration
	FormatPercent
	FormatOrdinal
	FormatComma
	FormatMetricPrefix
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

type (
	// BotVariable is what is used for the ActionConsideration scoring process.
	//
	// BotVariable is assigned in the BotGroup, which is the definition of what will be queried or synthesized into
	// each Bot.
	//
	// If Evaluate is not empty, then this will not run a Query, and instead will execute after all Query Variables are
	// set, and will Evalutate using govaluate.Evaluate() to set a new variable.
	//
	// Otherwise, a query is performed, and the variable is set from the query.
	//
	// Query Variables use any combination of BotKey, QueryKey and QueryKeyValue to set the variables.
	//
	// If BotKey is set, only query results that have a Metric Key named BotKey that matches Bot.Name will be accepted.
	//
	// If QueryKey is set, only query results that have a value with their Metric Name of QueryKey which matches
	// QueryKeyValue will be set.
	BotVariable struct {
		Type           BotVariableType   `json:"type"`
		Name           string            `json:"name"`
		Format         BotVariableFormat `json:"format"`
		BotKey         string            `json:"bot_key"` // Determines which Metric matches a Bot, may change between queries
		QueryName      string            `json:"query_name"`
		QueryKey       string            `json:"query_key"`       // Metric key to extract
		QueryKeyValue  string            `json:"query_key_value"` // Metric key value to match against the QueryKey
		Evaluate       string            `json:"evaluate"`        // If this is non-empty, query will not be performed.  After query testing for other variables, this will have a final phase of processing, and will take all the query-made variables and perform govaluation.Evaluate() with this evaluate string, to set this variable.  Evaluate variables cannot use each other, only Query variables.
		BoolRangeStart float64           `json:"bool_range_start"`
		BoolRangeEnd   float64           `json:"bool_range_end"`
		BoolInvert     bool              `json:"bool_invert"`
		Export         bool              `json:"export"` // If true, this variable will be exported for Metric collection.  Normally not useful, because we just got it from the Metric system.
	}
)

type (
	// BotGroup is used to create Bots.  Bots are the core of Sireus.  BotGroups define all the information used to
	// populate the ephemeral Bot structure.
	BotGroup struct {
		Name                   string                    `json:"name"`
		Info                   string                    `json:"info"`
		BotExtractor           BotExtractorQueryKey      `json:"bot_extractor"`             // This is the information we use to create the ephemeral Bots, but taking their names from this query's metric key
		States                 []BotForwardSequenceState `json:"states"`                    // States can only advance from the start to the end, they can never go backwards.  It's a sequence, but you can skip steps forward.  Using several of these, many situations can be modelled.
		LockTimers             []BotLockTimer            `json:"lock_timers"`               // Lock timers work at BotGroup or Bot level, and block any execution for a period of time, so the previous action's results can be evaluated
		BotTimeoutStale        util.Duration             `json:"bot_timeout_stale"`         // Duration since Bot.VariableValues was last updated until this Bot is marked as Stale.  Stale bots only execute Actions from a State named "Stale", so that you can respond, but no other states actions will apply.
		BotTimeoutRemove       util.Duration             `json:"bot_timeout_remove"`        // Duration since Bot.VariableValues was last updated until this bot is removed.  Bots are ephemeral.
		BotRemoveStoreDuration util.Duration             `json:"bot_remove_store_duration"` // Duration since removal that a Bot is stored for inspection, so that you don't lose access to useful information.  If the bot returns before this duration is over, it will be resumed.  Resumption can be refused setting BotGroup.RefuseBotResumption
		RefuseBotResumption    bool                      `json:"refuse_bot_resumption"`     // If true, once a bot is removed, while it is being stored for inspect, if it returns it will not be resumed.  Instead a new bot will be created to disconnect their history, even though they share the same BotKey
		ActionThreshold        float64                   `json:"action_threshold"`          // Minimum Action Final Score to execute a command.  Allows ignoring lower scoring Actions for testing or troubleshooting
		CommandHistoryDuration util.Duration             `json:"command_history_duration"`  // How long we keep history for ActionCommandResult values
		JournalRollupStates    []string                  `json:"journal_rollup_states"`     // If any of these states become Active, then they will be rolled up into Journal collection, example: Outage Report
		JournalRollupDuration  util.Duration             `json:"journal_rollup_duration"`   // Time between a Journal Rollup ending, and another Journal Rollup beginning, so that they are grouped together.  This collects flapping outages together.
		Queries                []BotQuery                `json:"queries"`                   // Queries used to populate the Variables
		Variables              []BotVariable             `json:"variables"`                 // Variables get their data from Queries, and are used in ActionConsideration evaluations to score the Action
		Actions                []Action                  `json:"actions"`                   // Actions get scored using ActionConsideration and the highest scored Action that IsAvailable will be executed.  Excecution also requires no LockTimers or other blocking factors.  The biggest factor is that Actions only are tested and execute when certain BotStates are set, so there is a built-in grouping of available Actions based on the BotState.
		Bots                   []Bot                     // These are the ephemeral workers of Sireus.  In an Action, the Queries populate VariableValues and then the ActionConsiderations are scored to determine if an action IsAvailable.

		// Invalid = Isn't getting all the information.  Stale = Information out of data.  Removed = No data for too long, removing.
		InvalidBots   []string
		StaleBots     []string
		RemovedBots   []string
		FreezeActions bool // If true, no actions will be taken for this BotGroup.  Allows group level control.
	}
)

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

type (
	// QueryServer is where we connect to get data to populate our Bots.  example: Prometheus
	// These are stored at a Site level, so that they can be shared by all BotGroups in a Site.
	//
	// Inside a QueryServer, all QueryNames must be unique for any BotGroup, so that they can potentially be shared
	// to reduce QueryServer traffic.  Keep this in mind when creating BotGroup.Queries.
	QueryServer struct {
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
)

type (
	// Top Level of the data structure.  Site silos all BotGroups and QueryServers, so that we can have multiple Sites
	// which are using different data sets, and should not share any data with each other.
	Site struct {
		Name          string        `json:"name"`
		Info          string        `json:"info"`
		BotGroupPaths []string      `json:"bot_group_paths"`
		QueryServers  []QueryServer `json:"query_servers"`
		BotGroups     []BotGroup
		FreezeActions bool // If true, no actions will be taken for this Site.  Allows control of all BotGroups.
	}
)
