# sireus
<img align="right" width="30%" src="docs/images/sireus_logo.png">

## Sireus - SRE Utility System - Dynamic Bot manager for executing SRE and DevOps commands conditionally

**NOTE: Sireus is in RFC stage and not suitable for production usage.  [Please file questions, comments and requests here.](https://github.com/ghowland/sireus/issues)**

Replaces cron jobs, Jenkins, Nagios, or other less sophisticated execution methods for SRE/DevOps automation.  

Sireus is a Decision System, made to collect information from Monitoring or other data sources, and make a decision on which action to execute, if any action should be executed.

<img width="70%" src="docs/images/sireus_stack_pos_exec.png" alt="Sireus Stack Position">

### Table of Contents

- [Sireus Goals](#sireus-goals)
- [Links](#links)
- [Data Structure](#data-structure)
  * [Sireus Bots and Bot Groups](#sireus-bots-and-bot-groups)
  * [Terminology](#terminology)
  * [How a Utility System or "Utility AI" works](#how-a-utility-system-or--utility-ai--works)
    + [Action Consideration Data](#action-consideration-data)
    + [Action Final Scores from Multiple Considerations](#action-final-scores-from-multiple-considerations)
    + [Why so many steps to get to a Final Action Score?](#why-so-many-steps-to-get-to-a-final-action-score-)
- [Sireus Portrait](#sireus-portrait)

### Sireus Goals

- Bots execute a *single* command or API call out of many possibilities; designed for SRE and DevOps environments.
- Sireus is a Decision System.  It's purpose is to make a decision and execute a *single* command or web call.
- Fits into the stack between monitoring and alerting.  ex: Prometheus -> Sireus -> Alert Manager.
- Works with existing software stack, with minimal configuration.  Architecture agnostic.
- Dynamically create Bots for any Platform, Service, Process, Host, etc from monitoring software (ex: Prometheus).  Bots are ephemeral.
- Bots have something like rulesets for prioritizing conditional commands to respond to detected issues.
- Scalable to large amounts of tests and commands, with deterministic execution, and inspectable with historical or test data to aid in configuration and adjusting values to better respond to future events.
- Locking commands per Bot or Bot Group, to stop conflicting commands from running at once, or within a window to verify results of previous commands.
- Uses the ["Utility AI" or "Utility System" behavior system](https://en.wikipedia.org/wiki/Utility_system), which provides a sophisticated method scoring for N conditions per command, to prioritize execution based on collected Bot information.

## Links

 - [Sireus Data Structure and Internal Function Documentation](docs/godoc.md)

## Data Structure

![Data Structure](/docs/images/data_structures.png)

### Sireus Bots and Bot Groups

- A Bot Group is defined statically to create Bots.  Queries against monitoring software (ex: Prometheus) or services (ex: Kubernetes) are defined in the Bot Group to be used by Bots.
- Bots are suggested to be created dynamically from monitoring data
- Bots can also be created statically, for less dynamic services (ex: Kafka)
- Bot Groups and Bots have arbitrary variables set with timeouts to ensure execution doesn't occur from stale data
- Triggers to execute commands for common functions, such as a Bot's data disappearing from monitoring data (stale or missing)
- Commands are meant to execute against a service or web API, host (ex: bash), or to update internal Sireus data for more complex conditional testing.  This allows building up more complex state variables, which are easier to read and reason about in the conditional logic.

### Terminology

- **Bot Group**: A collection of Bots, for executing Actions, based on conditional scoring.  This would be mapped against a Web App or other software service in your infrastructure.
- **Bot**: A collection of Variable Data and Actions, which contain conditional scoring information based on monitoring queries, which then executes a command.  Each Bot keeps information to use in making decisions.
- **Action**: This is the wrapper for conditions to create a Score, and the Command to execute if it is selected.
- **Action Score**: This is the priority of execution.  Given a set of potential Actions, we rank them from highest to lowest score, executing the highest score, and never execute Actions with a score of 0.
- **Action Consideration**: These are essentially conditions, but are floats to provide a range of data, instead of only boolean.
- **Action Command**: Executing 1 or more bash-type OS level commands or a service or web API calls.  Generalizing all of these to an "Action Command".

### How a Utility System or "Utility AI" works

- All configuration is defined per Bot Group.  These consist of a set of Actions.
- Each Action has a set of Considerations (Conditions that are not just boolean) which create a Score.
- The highest non-zero score will be executed.  In most cases, nothing will be done and all scores will be zero, because no actions are necessary.  When actions become necessary, the highest non-zero scored Action will be executed.

<img width="70%" src="docs/images/bot_action_execution.png" alt="Bot Action Execution">

#### Action Consideration Data

An Action has N Considerations, made from the following data:

- **Weight**: Per-consideration weight, so each consideration can have higher or lower weight than others
- **Value Function**: A function or command to execute to get a value (float)
- **Value Range**: A range of data ranges to test the result of the consideration's function output.  ex: 0.0-1.0, 0-100, 35-999.  This is the Floor and the Ceiling of the Value Function output.
- **Curve**: A curve to apply Value Function output.  The 2D Curve data goes from 0-1 on X and Y axis.  X is the Value Function Range position, and Y will be multiplied by the Weight to give the final Score.

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

In the above single Consideration Data, we had a single Consideration Score of 3.55.  If there were more considerations, all of these would be calculated together, to get a final consideration score, and then multiplied by the Action Weight to get a final Action Score.

**Example of an Action with Multiple Considerations:**

- **Action**: Send API Remediation XYZ
- **Action Weight**: 1.5
- **Final Calculated Scores for all Considerations**: 3.55
- **Final Action Score**: 5.32

When all the Actions have had their Final Scores calculated, if 5.32 is the highest score, then that action will be executed.  

For a given Action, if **any** of the Considerations have a score of zero, then the entire Final Action Score is zero.  This allows any Consideration to make an Action invalid.

#### Why so many steps to get to a Final Action Score?

The reason to have all of these steps is to be able to control exactly how important any given consideration test is to executing that action, and to provide multiple ways to invalidate the action (any consideration with a 0 score).

The benefit of this is that even with hundreds or thousands of Actions, they can be tuned so that the correct action executes at the correct time.  These tests are deterministic, and can be run on historic or test data, so that execution can be tested on prior outages to see how the rules would execute in known failure situations, or proposed failure situations using test data.

Having the ability to tune values at the top level Action, and for each Consideration, allows for a lot of tuning ability to ensure correct execution.

## Sireus Portrait

![Sireus Portrait](/docs/images/sireus_dog_star.png)
