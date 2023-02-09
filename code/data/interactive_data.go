package data

import "time"

type (
	// Pool to keep our all InteractiveSession data
	InteractiveSessionPool struct {
		Sessions []InteractiveSession // All our current InteractiveSession data, for tracking users testing scoring or config changes through the web app.  Will store an addition set of BotQuery items per BotGroup overridden
	}
)

type (
	// An InteractiveSession is created when a Web App user wants to look at how their Actions would score at a previous time, or if there were different Bot.VariableValues or an Action.Weight or ActionConsideration was different
	InteractiveSession struct {
		UUID          int64     `json:"uuid"`           // This is the unique identifier for this InteractiveSession, and cannot be 0.  0 is used by the normal server processes for performing queries.
		QueryTime     time.Time `json:"query_time"`     // Time to make all our queries, so that we can interactively look into past data and reply how actions would be scored with the current config (base and OverrideData)
		Override      Override  `json:"overrides"`      // This is a collection of data we get from the Web Client that overrides internal or queried data.  Over
		TimeRequested time.Time `json:"time_requested"` // This is the last time we received a request from this InteractiveSession.  When it passes the AppConfig.InteractiveSessionTimeout duration it will be removed
	}
)

type (
	// This tracks all the override changes relating to BotGroups or Bots for an InteractiveSession
	Override struct {
		BotGroups []OverrideBotGroup `json:"bot_groups"` // Overrides of BotGroup data: Action.Weight and ActionConsideration data
		Bots      []OverrideBot      `json:"bots"`       // Overrides of Bot data: Bot.VariableValues
	}
)

type (
	// Overrides to a BotGroup for an InteractiveSession
	OverrideBotGroup struct {
		BotGroupName         string                                 `json:"bot_group_name"`        // Name of the Bot to override.  This scope is Bot level
		ActionWeight         map[string]float64                     `json:"action_weight"`         // Overrides an Action.Weight for all the Bots in this BotGroup.  Changes Action scores for all Bots in a BotGroup
		ActionConsiderations map[string]OverrideActionConsideration `json:"action_considerations"` // Overrides ActionConsideration values for all Bots in this BotGroup
	}
)

type (
	// Overrides for an ActionConsideration in an Action, for a BotGroup.  Changes all related Bot scores
	OverrideActionConsideration struct {
	}
)

type (
	// Overrides to a Bot for an InteractiveSession
	OverrideBot struct {
		BotName           string             `json:"bot_name"`            // Name of the Bot to override.  This scope is Bot level
		BotVariableValues map[string]float64 `json:"bot_variable_values"` // The Bot.VariableValues that are being overridden.  This is useful to see how Action scores would change if some monitoring data was different, without having to find a time in the past where that was true.  Allows planning for different situations.
	}
)
