package data

import "time"

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
		Name                 string                   // Unique identifier pulled from the BotGroup.BotExtractor
		VariableValues       map[string]float64       // These are the unique values for this Bot, and will be used for all ActionConsideration scoring
		SortedVariableValues PairFloat64List          // Sorted VariableValues, Handlebars helper
		StateValues          []string                 // These are the current States for this Bot.  Actions can only be available for execution, if all their Action.RequiredStates are active in the Bot
		CommandHistory       []ActionCommandResult    // Storage of previous ActionCommand data run, so we can see insight into the history
		LockTimers           []BotLockTimer           // LockTimers allow control over Actions that require them, so they cant be available until they can get all their LockTimers
		ActionData           map[string]BotActionData // Key is Action.Name
		SortedActionData     PairBotActionDataList    // Scored ActionData, Handlebars helper
		FreezeActions        bool                     // If true, no actions will be taken for this Bot.  Single agent control
		IsInvalid            bool                     // If true, this Bot is Invalid and cannot make actions, because not all the Variables were found
		InfoInvalid          string                   // Short sentences ending with ".  " concatenated into this string to give all the reasons this Bot.IsInvalid
		IsStale              bool                     // If true, this Bot is Stale, and cannot make decisions.  IsInvalid is the super-state, and will be marked from this sub-reason for invalidity
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
