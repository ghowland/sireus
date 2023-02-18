package data

import (
	"time"
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
		BotTimeoutStale        Duration                  `json:"bot_timeout_stale"`         // Duration since Bot.VariableValues was last updated until this Bot is marked as Stale.  Stale bots only execute Conditions from a State named "Stale", so that you can respond, but no other states actions will apply.
		BotTimeoutRemove       Duration                  `json:"bot_timeout_remove"`        // Duration since Bot.VariableValues was last updated until this bot is removed.  Bots are ephemeral.
		BotRemoveStoreDuration Duration                  `json:"bot_remove_store_duration"` // Duration since removal that a Bot is stored for inspection, so that you don't lose access to useful information.  If the bot returns before this duration is over, it will be resumed.  Resumption can be refused setting BotGroup.RefuseBotResumption
		RefuseBotResumption    bool                      `json:"refuse_bot_resumption"`     // If true, once a bot is removed, while it is being stored for inspect, if it returns it will not be resumed.  Instead a new bot will be created to disconnect their history, even though they share the same BotKey
		ConditionThreshold     float64                   `json:"action_threshold"`          // Minimum Condition Final Score to execute a command.  Allows ignoring lower scoring Conditions for testing or troubleshooting
		CommandHistoryDuration Duration                  `json:"command_history_duration"`  // How long we keep history for ConditionCommandResult values
		JournalRollupStates    []string                  `json:"journal_rollup_states"`     // If any of these states become Active, then they will be rolled up into Journal collection, example: Outage Report
		JournalRollupDuration  Duration                  `json:"journal_rollup_duration"`   // Time between a Journal Rollup ending, and another Journal Rollup beginning, so that they are grouped together.  This collects flapping outages together.
		Queries                []BotQuery                `json:"queries"`                   // Queries used to populate the Variables
		Variables              []BotVariable             `json:"variables"`                 // Variables get their data from Queries, and are used in ConditionConsideration evaluations to score the Condition
		Conditions             []Condition               `json:"actions"`                   // Conditions get scored using ConditionConsideration and the highest scored Condition that IsAvailable will be executed.  Excecution also requires no LockTimers or other blocking factors.  The biggest factor is that Conditions only are tested and execute when certain BotStates are set, so there is a built-in grouping of available Conditions based on the BotState.
		Bots                   []Bot                     // These are the ephemeral workers of Sireus.  In a Condition, the Queries populate VariableValues and then the ConditionConsiderations are scored to determine if an action IsAvailable.

		// Invalid = Isn't getting all the information.  Stale = Information out of data.  Removed = No data for too long, removing.
		InvalidBots      []string
		StaleBots        []string
		RemovedBots      []string
		FreezeConditions bool   // If true, no actions will be taken for this BotGroup.  Allows group level control.
		LockKey          string // Formatted with: (Site.Name).(BotGroup.Name)
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

type (
	// Forward Sequence State is the term I am using to describe a State Machine that only has a single forward
	// sequence.  It can be Advanced and it can be Reset, but the state cannot go backwards.
	//
	// In this way you can create a State Machine for investigating problems, trying to solve them, checking for
	// resolution, and finally escalating and waiting for someone else to fix it.
	//
	// If a resolution is detected by a Condition, the action can Reset this state, starting the States process over again.
	//
	// States are used to exclude Conditions from being tested, so that Conditions can be targetted at a specific State of a
	// Bot's operation.  This allows segmenting logic.  Conditions use Condition.RequiredStates to limit when they can execute.
	BotForwardSequenceState struct {
		Name   string   `json:"name"`
		Info   string   `json:"info"`
		Labels []string `json:"labels"`
	}
)

type (
	// Differentiate Query Types, so we can format our requests and parse the data properly
	BotQueryType int64
)

const (
	Range BotQueryType = iota
)

// Format the BotQueryType to a string usable for building the request
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
		QueryServer string       `json:"query_server"`
		QueryType   BotQueryType `json:"query_type"`
		Name        string       `json:"name"`
		Info        string       `json:"info"`
		Query       string       `json:"query"`
		Interval    Duration     `json:"interval"`
	}
)

type (
	// Scope for locking Conditions
	BotLockTimerType int64
)

const (
	LockBotGroup BotLockTimerType = iota
	LockBot
)

// Format the BotLockTimerType for human readability
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
	// BotLockTimer is used to both block a Condition from being executed, if the BotLockTimer.IsActive and has not
	// reached the Timeout yet.  Conditions can use multiple BotLockTimers which essentially act as execution "channels"
	// where Conditions execute 1 at a time.
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
		ActivatedByBot string // Bot.Name of who set this Lock Timer, so we can track Conditions
	}
)

type (
	// This is the raw input data type.  It will still be turned into a float64, but it is best to know the origin type
	BotVariableType int64
)

const (
	Boolean BotVariableType = iota
	Float
)

type (
	// This is for formatting the data we got raw from BotVariableType.  This uses Humanize and other readability funcs
	BotVariableFormat int64
)

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

// Format the BotVariableType for human readability
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
	// BotVariable is what is used for the ConditionConsideration scoring process.
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
