package data

import "time"

type (
	// State Condition is what is considered for execution.  It will receive a Final Score from its Weight and Consideration Final Scores
	Condition struct {
		Name               string                   `json:"name"`                 // Name of the Condition
		Info               string                   `json:"info"`                 // Description
		IsLaunched         bool                     `json:"is_launched"`          // If false, this will never execute.  Launching means it is configured and ready to run live.  When Conditions are created, is_launched==false and must be changed so that the action could execute.
		IsDisabled         bool                     `json:"is_disabled"`          // When testing changes, disable with modifying config
		Weight             float64                  `json:"weight"`               // This is the multiplier for the Final Score, from the Consideration Final Score
		WeightMin          float64                  `json:"weight_min"`           // If Weight != 0, then this is the Floor value.  We will bump it to this value, if it is less than this value
		WeightThreshold    float64                  `json:"weight_threshold"`     // If non-0, this is the threshold to be Active, and potentially execute Conditions.  If the Final Score is less than this Threshold, this Condition can never run.  WeightMin and WeightThreshold are independent tests, and will have different results when used together, so take that into consideration.
		ExecuteRepeatDelay Duration                 `json:"execute_repeat_delay"` // Duration until this Condition can execute again.  If short, this just the problem of double execution if it is 0, which is required.  It can't be 0.  If this is long, this becomes a good way to process other actions instead of this one, because you already tried it recently.
		RequiredAvailable  Duration                 `json:"required_available"`   // If greater than 0s, this Condition must have been continuously Available for this Duration for it to be executed.  Allows us to make sure it's not flapping or inconsistent for a period of time before being executed
		RequiredLockTimers []string                 `json:"required_lock_timers"` // All of these Lock Timers must be available for this Condition to trigger.  Afterwards, they will all be locked for ConditionCommand.LockTimerDuration automatically
		RequiredStates     []string                 `json:"required_states"`      // All of these states must be Active for this
		Considerations     []ConditionConsideration `json:"considerations"`       // These Considerations are used to create a Score for this Condition, which must be the highest score, and must be higher than the MinimumThreshold, and if all other requirements are met, this Condition will be executed
		Command            ConditionCommand         `json:"command"`              // This is the command that will be executed.  It could just change States, or run a Command or API call
	}
)

type (
	// Considerations are units for scoring a State Condition.  Each creates a Score, and they are combined to create the
	// Consideration Final Score.
	ConditionConsideration struct {
		Name       string  `json:"name"`
		Weight     float64 `json:"weight"`
		CurveName  string  `json:"curve"`
		RangeStart float64 `json:"range_start"`
		RangeEnd   float64 `json:"range_end"`
		Evaluate   string  `json:"evaluate"`
	}
)

type (
	// What will we do with this ConditionCommand?  We will only ever do 1 thing per Condition, as this is a Decision System
	ConditionCommandType int64
)

const (
	ShellCommand    ConditionCommandType = iota // Shell Command
	WebHttps                                    // HTTPS request
	WebHttpInsecure                             // HTTP request
	WebRPC                                      // RPC call
	NoOperation                                 // Do nothing
)

func (act ConditionCommandType) String() string {
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
	// When a Condition is selected for execution by its Final Score, the ConditionCommand is executed.  A command or web request.
	ConditionCommand struct {
		Name              string               `json:"name"`                // Best Practice: Description of what this command is going to do.  Focus on the effect this will cause, and what it affects.
		LogFormat         string               `json:"log_format"`          // This is what will be logged for human readability for the Name of this command.  It is formatted by Handlebars and can access data from: bot, botGroup, action, actionCommand.  Best practices, expand into a specific target for the Name field's general purpose.
		Type              ConditionCommandType `json:"type"`                // Type of command that was executed
		Content           string               `json:"content"`             // Payload of the command or RPC
		SuccessStatus     int                  `json:"success_status"`      // Success or failure?
		SuccessContent    string               `json:"success_content"`     // Data received from endpoint about our SuccessState and any return payload
		LockTimerDuration Duration             `json:"lock_timer_duration"` // We will set the Condition.RequiredLockTimers to this duration to block anything from running.  Each of them all had to be available, and now will all be blocked.  Design your structures around this concept.  LockTimers provide "lanes" of execution that can overlap or work independently.
		HostExecKey       string               `json:"host_exec_key"`       // Sireus Client presents this key to get commands to run
		SetBotStates      []string             `json:"set_bot_states"`      // Will Advance all of these Bot States.  Advance can only go forward in the list, or start at the very beginning.  It can't go backwards, that is invalid data.  Only the State.Name and not the StateName.State is present, this will just advance to the next available state until it hits the final one and stay there.
		ResetBotStates    []string             `json:"reset_bot_states"`    // Will reset all these Bot States to their first entry.  This is how Sireus handles state flow: forward-only and then reset
	}
)

type (
	// When a Condition is selected for execution by its Final Score, the ConditionCommand will execute and store this result.
	ConditionCommandResult struct {
		BotGroupName  string    // Name of the BotGroup that had this Bot Condition, for easy lookup
		BotName       string    // Name of the Bot that had this condition, for easy lookup
		ConditionName string    // Name of the Condition that had this command, for easy lookup
		CommandLog    string    // ConditionCommand.LogFormat gets formatted and put here for rich markup over Name field
		ResultStatus  string    // Status of the result (success/failure)
		ResultContent string    // Content of the result, for full inspection
		HostExecOn    string    // Host this command was executed on, given by Sireus Client
		Started       time.Time // Time this command started
		Finished      time.Time // Time this command finishing
		Score         float64   // This was the Condition Final Score
		StatesBefore  []string  // These are the Bot.StateValues before this command was run
		StatesAfter   []string  // These are the Bot.StateValues after this command was run
	}
)
