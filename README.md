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

![Bot Group](/docs/images/bot_group.png)

#### Action Consideration Data

An Action has N Considerations, made from the following data:

- **Weight**: Per-consideration weight, so each consideration can have higher or lower weight than others
- **Value Function**: A function or command to execute to get a value (float)
- **Value Range**: A range of data ranges to test the result of the consideration's function output.  ex: 0.0-1.0, 0-100, 35-999.  This is the Floor and the Ceiling of the Value Function output.
- **Curve**: A curve to apply Value Function output.  The 2D Curve data goes from 0-1 on X and Y axis.  X is the Value Function Range position, and Y will multipled by the Weight to give the final Score.

**Example a Single Consideration:**

- Weight: 5.0
- Value Function Result: 60
- Value Range: 0 to 100
- Curve:

![Curve Example](/docs/images/curve_example.PNG)

Given a Value Function Result (60) in the Value Range (0 to 100) = 0.6

In the Curve, with the X=0.6 the Y value = 0.71

The Curve Result (0.71) is multiplied by the Weight (5): 0.71 * 5 = 3.55 Consideration Score

#### Action Final Scores from Multiple Considerations

In the above single Consideration Data, we had a single Consideration Score of 3.55.  If there were more considerations, all of these would be calculated together, to get a final consideration score, and then multipled by the Action Weight to get a final Action Score.

**Example of an Action with Multiple Considerations:**

- **Action**: Send API Remediation XYZ
- **Action Weight**: 1.5
- **Final Calculated Scores for all Considerations**: 3.55
- **Final Action Score**: 5.32

When all the Actions have had their Final Scores calculated, if 5.32 is the highest score, then that action will be executed.  

For a given Action, if **any** of the Considerations have a score of zero, then the entire Final Action Score is zero.  This allows any Consideration to make an Aciton invalid.

#### Why so many steps to get to a Final Action Score?

The reason to have all of these steps is to be able to control exactly how important any given consideration test is to executing that action, and to provide multiple ways to invalidate the action (any consideration with a 0 score).

The benefit of this is that even with hundreds or thousands of Actions, they can be tuned so that the correct action executes at the correct time.  These tests are deterministic, and can be run on historic or test data, so that execution can be tested on prior outages to see how the rules would execute in known failure situations, or proposed failure situations using test data.

Having the ability to tune values at the top level Action, and for each Consideration, allows for a lot of tuning ability to ensure correct execution.
