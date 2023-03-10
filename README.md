# sireus
<img align="right" alt="Sireus Logo" width="30%" src="https://github.com/ghowland/sireus/blob/main/documentation/images/sireus_logo.png">

## Sireus - SRE Utility System
### Decision System for tracking SRE and DevOps operational state and executing commands

**NOTE: Sireus is in RFC stage and not suitable for production usage.  [Please file questions, comments and requests here.](https://github.com/ghowland/sireus/issues)**

Replaces cron jobs, Jenkins, Nagios, or other less sophisticated execution methods for SRE/DevOps automation.  

Sireus is a Decision System, made to collect information from Monitoring or other data sources, and make a decision on which condition command to execute, if any condition should be executed.

<img width="70%" src="https://github.com/ghowland/sireus/blob/main/documentation/images/sireus_stack_pos_exec.png" alt="Sireus Stack Position">

### Table of Contents

- [Sireus Goals](#sireus-goals)
- [Links to Documentation and Communication Options](#links)
- [Data Structure](#data-structure)
  * [Sireus Bots and Bot Groups](#sireus-bots-and-bot-groups)
  * [Terminology](#terminology)
  * [How a Utility System or "Utility AI" works](#how-a-utility-system-or-utility-ai-works)
    + [State Condition Consideration Data](#state-condition-consideration-data)
    + [State Condition Final Scores from Multiple Considerations](#state-condition-final-scores-from-multiple-considerations)
    + [Why so many steps to get to a Final State Condition Score?](#why-so-many-steps-to-get-to-a-final-state-condition-score)
- [Help Wanted... in many areas including Data Visualization and Web Design](#help-wanted)
- [Planned Features](#planned-features)
- [Sireus Portrait](#sireus-portrait)

### Sireus Goals

- Bots execute a *single* command or API call out of many possibilities; designed for SRE and DevOps environments.
- Sireus is a Decision System.  Its purpose is to make a decision and execute a *single* command or web call.
- Fits into the stack between monitoring and alerting.  ex: Prometheus -> Sireus -> Alert Manager.
- Works with existing software stack, with minimal configuration.  Architecture agnostic.
- Dynamically create Bots for any Platform, Service, Process, Host, etc. from monitoring software (ex: Prometheus).  Bots are ephemeral.
- Bots have something like rule sets for prioritizing conditional commands to respond to detected issues.
- Scalable to large amounts of tests and commands, with deterministic execution, and inspectable with historical or test data to aid in configuration and adjusting values to better respond to future events.
- Locking commands per Bot or Bot Group, to stop conflicting commands from running at once, or within a window to verify results of previous commands.
- Uses the ["Utility AI" or "Utility System" behavior system](https://en.wikipedia.org/wiki/Utility_system), which provides a sophisticated method scoring for N conditions per command, to prioritize execution based on collected Bot information.  Scales to large numbers of commands, allowing for complex reactions in large environments.
- Limited scope.  A key goal of Sireus is a limited scope, so it remains a focused tool.  The goal is to run commands based on States and Conditions that are met, and then log that and make all the information visible.  There is some depth and complexity, but there is also a natural stopping point where Sireus should never go beyond, so it becomes stable and trustworthy in it's original mission.


## Links

 - [Blog](https://blog.sireus.cloud/)
 - [How to Start Configuring Sireus in 10 Steps](https://github.com/ghowland/sireus/blob/main/documentation/how_to_start.md)
 - [Best Practices](https://github.com/ghowland/sireus/blob/main/documentation/best_practices.md)
 - [Discord](https://discord.gg/VTVXrXJWxk)
 - [Data Structure and Internal Function Documentation](https://github.com/ghowland/sireus/blob/main/documentation/godoc.md)
 - [Contributing](https://github.com/ghowland/sireus/blob/main/documentation/contributing.md)
 - [Developer Chat on Zulip - Invite only for now](https://sireus.zulipchat.com/) - If you want to join the development process, please start by [creating Issues.](https://github.com/ghowland/sireus/issues)  Sireus is currently in the Design RFC phase.
 - How to pronounce Sireus?  Like the word "serious".
 - Web App Example Page:

![Web App Example Page](https://github.com/ghowland/sireus/blob/main/documentation/images/webapp_example.png)

### Sireus Bots and Bot Groups

- A Bot Group is defined statically to create Bots.  Queries against monitoring software (ex: Prometheus) or services (ex: Kubernetes) are defined in the Bot Group to be used by Bots.
- Bots are suggested to be created dynamically from monitoring data
- Bots can also be created statically, for less dynamic services (ex: Kafka)
- Bot Groups and Bots have arbitrary variables set with timeouts to ensure execution doesn't occur from stale data
- Triggers to execute commands for common functions, such as a Bots data disappearing from monitoring data (stale or missing)
- Commands are meant to execute against a service or web API, host (ex: bash), or to update internal Sireus data for more complex conditional testing.  This allows building up more complex state variables, which are easier to read and reason about in the conditional logic.

## Data Structure

![Data Structure](https://github.com/ghowland/sireus/blob/main/documentation/images/data_structure.png)

### Terminology

The data structure image above shows the relation between these terms:

- **Bot Group**: A collection of Bots, for executing State Conditions, based on conditional scoring.  This would be mapped against a Web App or other software service in your infrastructure.
- **Bot**: A collection of Variable Data and State Conditions, which contain conditional scoring information based on monitoring queries, which then executes a command.  Each Bot keeps information to use in making decisions.
- **State**: Bot Groups describe all the states a Bot can be in in the form of Forward Sequence State pipelines.  This gives several "pipelines" of states that only move or skip forward or reset to the beginning, so they never get stuck in an inner loop.  The current set of States a Bot has will determine was State Conditions are available, and ultimately what State Condition Commands can be executed.
- **State Condition**: This is the wrapper for conditions to create a Score, and the Command to execute if it is selected.
- **State Condition Score**: This is the priority of execution.  Given a set of potential State Conditions, we rank them from highest to lowest score, executing the highest score, and never execute State Conditions with a score of 0.
- **State Condition Consideration**: These are essentially conditions, but are floats to provide a range of data, instead of only boolean.
- **State Condition Command**: Executing a single bash-type OS level commands or a service or web API calls.

### How a Utility System or Utility AI works

- All configuration is defined per Bot Group.  These consist of a set of State Conditions.
- Each State Condition has a set of Considerations (Conditions that are not just boolean) which create a Score.
- The highest non-zero score will be executed.  In most cases, nothing will be done and all scores will be zero, because no actions are necessary.  When actions become necessary, the highest non-zero scored State Condition will be executed.

<img width="70%" src="https://github.com/ghowland/sireus/blob/main/documentation/images/bot_condition_execution.png" alt="Bot State Condition Execution">

#### State Condition Consideration Data

A State Condition has N Considerations, made from the following data:

- **Weight**: Per-consideration weight, so each consideration can have higher or lower weight than others
- **Value Function**: A function or command to execute to get a value (float)
- **Value Range**: A range of data ranges to test the result of the consideration's function output.  ex: 0.0-1.0, 0-100, 35-999.  This is the Floor and the Ceiling of the Value Function output.
- **Curve**: A curve to apply Value Function output.  The 2D Curve data goes from 0-1 on X and Y axis.  X is the Value Function Range position, and Y will be multiplied by the Weight to give the final Score.

**Example a Single Consideration:**

- Weight: 5.0
- Value Function Result: 60
- Value Range: 0 to 100
- Curve:

![Curve Example](https://github.com/ghowland/sireus/blob/main/documentation/images/curve_example.PNG)

Given a Value Function Result (60) in the Value Range (0 to 100) = 0.6

In the Curve, with the X=0.6 the Y value = 0.71

The Curve Result (0.71) is multiplied by the Weight (5): 0.71 * 5 = 3.55 Consideration Score

#### State Condition Final Scores from Multiple Considerations

In the above single Consideration Data, we had a single Consideration Score of 3.55.  If there were more considerations, all of these would be calculated together, to get a final consideration score, and then multiplied by the State Condition Weight to get a final State Condition Score.

**Example of a State Condition with Multiple Considerations:**

- **State Condition**: Send API Remediation XYZ
- **State Condition Weight**: 1.5
- **Final Calculated Scores for all Considerations**: 3.55
- **Final State Condition Score**: 5.32

When all the State Conditions have had their Final Scores calculated, if 5.32 is the highest score, then that condition will be executed.  

For a given State Condition, if **any** of the Considerations have a score of zero, then the entire Final State Condition Score is zero.  This allows any Consideration to make a State Condition invalid.

#### Why so many steps to get to a Final State Condition Score?

The reason to have all of these steps is to be able to control exactly how important any given consideration test is to executing that condition, and to provide multiple ways to invalidate the condition (any consideration with a 0 score).

The benefit of this is that even with hundreds or thousands of State Conditions, they can be tuned so that the correct condition executes at the correct time.  These tests are deterministic, and can be run on historic or test data, so that execution can be tested on prior outages to see how the rules would execute in known failure situations, or proposed failure situations using test data.

Having the ability to tune values at the top level State Condition, and for each Consideration, allows for a lot of tuning ability to ensure correct execution.

### Help Wanted

A Decision System being used in SRE and DevOps land is a new tool, and how to represent and work with the data is not yet explored.  Looking for people to help make this easier for users to learn and become experts in, to create better automation and operational outcomes.

- Data visualization mysteries
	* How should the current state of a Bot Group be represented so that it can be understood at a glance?
		+ In the demo I (ghowland) show a list of the Bot Groups States, and how many Bots are in each state.  This gives some information, but I think much more information could be represented in a very brief manner and need someone to help figure this and other problems out.
		+ Another mystery to solve is how best to show the scoring values and curves.  My (ghowland) thought's on this are that there should be a simple-mode that is normally presented, which is just a boolean system, and hides the underlying scoring system, but still the same scoring, so it's unified.  Then an advanced system can open up all the scoring values as they are shown in the demo.  But, this needs to get designed.  I'll take a first pass at it soon.
- Web page design improvements for readability.
	* I (ghowland) did my best to keep it simple, but someone with an eye for design would be really helpful in making the pages easier to read and thus easier to gain insight from.
		+ All the web pages are rendered with [Handlebars](https://github.com/aymerick/raymond) (a [Mustache](https://github.com/cbroglie/mustache)-like), and I do almost all the processing using Handlebars Registered Helper system, where you would have full access to the data in the application, and then just use the handlebars syntax to loop over stuff or set a current context.  It's becoming a fairly robust library for this initial version's data representation.  It's also very easy to add any new Helpers, and my policy is to just add one for every condition as a 1-1 mapping of "I want to do X".  And of course reusing the existing ones as much as possible, but with an eye to not make any sneaky use cases, just a straight forward "Need to Verb with Adjective Noun" mappings.
- Development
	* A small plugin-system, so that custom functions could be called throughout the pipeline of the system.
		+ I (ghowland) think it's best to start this small with a minimal interface, and it can be kept as a legacy implementation when we find out what all the additional requirements we learn from use.  It needs a good first start to be a useful feedback tool so that users can spend the time to develop the expertise needed to push the system to its current limits, and get our feature set for the more mature plugin-system.  Because it's golang, I think it is best to just have them compile the plugins in, so can also avoid the expense of dynamic plugins.
- Organizing Documentation
	* I (ghowland) will create as much documentation as I think is needed to cover explaining the various areas and use cases, but I could use help organizing it.  Often as the primary author it's unclear what is confusing or understandable, or what information should be presented first.  I think 2 use cases of "As a new user" and "As an experienced user" could be tracked to present the information succinctly in any given page.  Is it for new users or experienced users?
		+ There are pretty good standards for this from other projects, so it would probably be best to just pick a successful project and model Sireus documentation's organization after theirs.  This decision hasn't been made yet for who to copy.
		+ A big help would just be suggestions or pull requests on re-organizing any existing documents to be easier to understand, please give the reason why it would make it easier to understand in the PR.
- Feedback
	* Sireus is in the Design RFC phase.  I want to get feedback on whether it is understandable, if not what are areas that lack clarity.  "What is it for?"  "Why should I use this?"  "How would I implement it?"  I have some of this information here now, but I (ghowland) don't know what is clear and what is unclear without more feedback.  [Please file questions, comments and requests here](https://github.com/ghowland/sireus/issues).

## Planned Features

These are a few planned features which are not in the current demo, but will be in the next version which is made to run in production sites:

- Authorization Sources: Per Site, like Query Servers.  These will call a shell command or URL to collect authorization information (bearer tokens, etc), which can be used in shell commands or URL calls for Condition Commands.  There will be a separate version of the values for URLs or commands printed for logging, so that this is designed to log safely.
- Sireus Client: An internal (goroutine) or stand alone client that will connect to the Sireus server, and present a set of "Host Keys" which determine what kind of shell commands or URLs can be run by this client.  This allows remote distribution of execution, and being able to target which jobs are run where.  Example: running a different shell command in AWS versus Google Cloud by Host Key "aws_host" or "gcloud_host" in the Condition Command config.
- Performance improvements: I implemented the initial prototype using maps in Go Lang, which require a lot of locking.  This delays serving pages, so I will switch all of those to slices to fix the performance issues when many queries are being made.

## Sireus Portrait

![Sireus Portrait](https://github.com/ghowland/sireus/blob/main/documentation/images/sireus_dog_star.png)
