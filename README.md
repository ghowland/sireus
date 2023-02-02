# sireus
## Sireus - SRE Utility System - Dynamic Bot manager for executing SRE and DevOps commands conditionally

Replaces cron jobs, Jenkins, Nagios, or other less sophisticated execution methods for SRE/DevOps automation.

### Sireus Goals

- Bots execute commands or API calls, designed for SRE and DevOps environments
- Dynamically create Bots for any Platform, Service, Process, Host, etc from monitoring software (ex: Prometheus)
- Bots have rulesets for prioritizing conditional commands to respond to detected issues
- Scalable to large amounts of tests and commands, with deterministic execution, and inspectable with historical or test data
- Locking commands per Bot or Bot Group, to stop conflicting commands from running at once, or within a window to verify results of previous commands.
- Uses the ["Utility AI" behavior system](https://en.wikipedia.org/wiki/Utility_system), which provides a sophisticated method scoring for N conditions per command, to prioritize execution based on collected Bot information 

### Sireus Bots and Bot Groups

- A Bot Group is defined statically to create Bots.  Queries against monitoring software (ex: Prometheus) or servies (ex: Kubernetes) are defined in the Bot Group to be used by Bots.
- Bots are suggested to be created dynamically from monitoring data
- Bots can also be created statically, for less dynamic services (ex: Kafka)
- Bot Groups and Bots have arbitrary variables set with timeouts to ensure the data is not stale
- Triggers to execute commands for common functions, such as a Bot's data disappearing from monitoring data (stale or missing)
- Commands are meant to execute against a service or web API, host (ex: bash), or to update internal Sireus data for more complex conditional testing.  This allows building up more complex state variables, which are easier to read and reason about in the conditional logic.

### Terminology

- **Bot Group**: A collection of Bots, for executing Actions, based on conditional scoring.
- **Bot**: A collection of Actions, which contain conditional scoring information based on monitoring queries, which then executes a command.
- **Action**: This is the wrapper for conditions to create a Score, and the Command to execute if it is selected
- **Score**: This is the priority of execution.  Given a set of potential commands, we rank them from highest to lowest score, executing the highest score, and never executing scores of 0.
- **Consideration**: These are essentially conditions, but can be a range of data, instead of only boolean.
- **Command**: Executing 1 or more: bash-type OS level or a service or web API.  Generalizing all of these to "Command"

### How a Utility System or "Utility AI" works

- All configuration is defined per Bot Group.  These consist of a set of Actions.
- Each Action has a set of Considerations (Conditions that are not just boolean) which create a Score.
- The highest non-zero score will be executed.  In most cases, nothing will be done and all scores will be zero, because no actions are necessary.  When actions become necessary, the highest non-zero scored Action will be executed.

