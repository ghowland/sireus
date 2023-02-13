package data

import (
	"sync"
	"time"
)

type (
	// Our session UUID.  UUID==0 is production data queried now.  UUID!=0 is interactive session with modified query and result data.
	SessionUUID int64
)

type (
	// Pool to keep our all InteractiveSession data.  This is used for normal production data too, with the UUID==0
	InteractiveSessionPool struct {
		Sessions   map[SessionUUID]InteractiveSession // All our current InteractiveSession data, key on UUID, for tracking users testing scoring or config changes through the web app.  Will store an addition set of BotQuery items per BotGroup overridden
		AccessLock sync.RWMutex                       // Lock for safe goroutine access to Sesssions map
	}
)

type (
	// An InteractiveSession is created when a Web App user wants to look at how their Actions would score at a previous time, or if there were different Bot.VariableValues or an Action.Weight or ActionConsideration was different
	InteractiveSession struct {
		UUID                     SessionUUID `json:"uuid"`             // This is the unique identifier for this InteractiveSession, and cannot be 0.  0 is used by the normal server processes for performing queries.
		TimeRequested            time.Time   `json:"time_requested"`   // This is the last time we received a request from this InteractiveSession.  When it passes the AppConfig.InteractiveSessionTimeout duration it will be removed
		QueryStartTime           time.Time   `json:"query_start_time"` // Time to make all our queries, so that we can interactively look into past data and reply how actions would be scored with the current config (base and OverrideData)
		QueryDuration            Duration    `json:"query_duration"`   // Duration to query past QueryStartTime
		QueryScrubTime           time.Time   `json:"query_scrub_time"` // Time our Scrubber is currently on.  Somewhere between QueryStartTime and QueryStartTime + QueryDuration
		Override                 Override    `json:"overrides"`        // This is a collection of data we get from the Web Client that overrides internal or queried data.  Over
		BotGroups                []BotGroup  // These are the BotGroups we create and cache for this Session
		IgnoreCacheQueryMismatch bool        // For Interactive sessions using the QueryStartTime and QueryDuration, this is true
		IgnoreCacheOverInterval  bool        // For Production sessions, this is true
	}
)

type (
	// This tracks all the override changes relating to BotGroups or Bots for an InteractiveSession
	Override struct {
		BotGroups []OverrideBotGroup `json:"bot_groups"` // Overrides of BotGroup data: Action.Weight and ActionConsideration data
		Bots      []OverrideBot      `json:"bots"`       // Overrides of Bot data: Bot.VariableValues and Bot.StateValues
		States    []string           `json:"states"`     // Overrides of States.  Can only set 1 per state.  By default, all should be set, so we know where we are.
	}
)

type (
	// Overrides to a BotGroup for an InteractiveSession
	OverrideBotGroup struct {
		BotGroupName         string                        `json:"name"`                  // Name of the Bot to override.  This scope is Bot level
		ActionWeight         map[string]float64            `json:"action_weight"`         // Overrides an Action.Weight for all the Bots in this BotGroup.  Changes Action scores for all Bots in a BotGroup
		ActionConsiderations []OverrideActionConsideration `json:"action_considerations"` // Overrides ActionConsideration values for all Bots in this BotGroup
	}
)

type (
	// Overrides for an ActionConsideration in an Action, for a BotGroup.  Changes all related Bot scores.  For simplicity, when making an ActionCconsideration override, all values are always updated.  No reason to have sparse changes here
	OverrideActionConsideration struct {
		ActionName        string  `json:"action_name"`        // Name of the Action to modify.  There are many ActionConsideration per Action
		ConsiderationName string  `json:"consideration_name"` // Consideration name identifier
		Weight            float64 `json:"weight"`             // Overrides ActionConsideration.Weight
		CurveName         string  `json:"curve_name"`         // Overrides ActionConsideration.CurveName
		RangeStart        float64 `json:"range_start"`        // Overrides ActionConsideration.RangeStart
		RangeEnd          float64 `json:"range_end"`          // Overrides ActionConsideration.RangeEnd
		//Evaluate string `json:"evaluate"` //TODO(ghowland): Not yet implemented.  This will be implemented later, leaving as reminder
	}
)

type (
	// Overrides to a Bot for an InteractiveSession
	OverrideBot struct {
		Name           string             `json:"name"`            // Name of the Bot to override.  This scope is Bot level
		VariableValues map[string]float64 `json:"variable_values"` // The Bot.VariableValues that are being overridden.  This is useful to see how Action scores would change if some monitoring data was different, without having to find a time in the past where that was true.  Allows planning for different situations.
		StateValues    []string           `json:"state_values"`    // If this is not an empty list, then it will override the specified BotGroup.States for this Bot.  This allows testing different Action scores in any state.
	}
)

type (
	// InteractiveControl comes from the Web App as JSON data on each RPC call, to update what data we want returned
	InteractiveControl struct {
		SessionUUID            SessionUUID `json:"sessionUUID"`            // Interactive Session UUID
		UseInteractiveSession  bool        `json:"useInteractiveSession"`  // If true, we will use all the Interactive session data.  If false, this will use normal production data.
		UseInteractiveOverride bool        `json:"useInteractiveOverride"` // If true, we will use the state and variable overrides.  If false, we can still use Query Time for interactive session, but will not use override states and vars
		PlayForward            bool        `json:"playForward"`            // Should we move the Scrubber forward as we return data?
		QueryStartTime         float64     `json:"queryStartTime"`         // When should the Interactive query start?
		QueryDuration          float64     `json:"queryDuration"`          // Duration from the Query Start, creates the Query End
		QueryScrubTime         float64     `json:"queryScrubTime"`         // This is where the Scrubbed currently is, so use this data from the metrics data
	}
)
