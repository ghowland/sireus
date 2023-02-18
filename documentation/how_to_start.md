# How to Start Configuring Sireus in 10 Steps

## 1. Determine Your Bot Groups

Map your Bot Groups to your services.  Anything you monitor already as a kind of specific target, like machines, pods, routers, software instances will all make good Bot Groups.  

Think of this like services you want to track, where you store a metric that identifies all the different instances of that server.  For example, if you had VMs and had a "host" metric (ex: "host=machinename"), then a Bot Group made of Bots from "host" would work well.

You can also use something like AWS or Google Cloud as the Bot Group, and different accounts as the Bots.  Ensure your metrics have a unique key to differentiate 1+ Bot in your Bot Group.  

Usually there are many Bots in a Bot Group, but sometimes you may be dealing with a service that is dealt with as a whole.  In this case you could just have 1 Bot in the Bot Group deal with the entire service's state.  This might be useful for vendors who have an API end point that you only interface with as a single entity.

## 2. For Every Bot Group, Pick a Bot Key to Extract Bots

Bots are made from a single Bot Group query, using a "Bot Key", which is just a Metric Name (ex: "host") which will allow the ephemeral creation of any number of Bots in that Bot Group.

Ensure all your desired metrics have this Boy Key value to match with the node.  The Metric Name can change per query, so if it is non-uniform that is OK, as long as the Metric Name values are uniform to match with the bot.  Ex: "host=instance1059"

## 3. For Each Bot Group, Create a Troubleshooting State Pipeline

For example, create a state named "Operation" and give it the following pipeline:

- Normal -> Problem -> Correction -> Evaluate -> Escalate -> EscalateWait

All states are a "forward sequence" so they only move from left to right, and then can be reset to the left again.  They can skips steps going to the right.  When the hit the right value, they stay there until they are reset.  It is a one-way state system, because this always leads to a conclusion.

We are not truly mapping the state of an application like a Finite State Machine, instead we are coming up with a pipeline for how to deal with data.  Different state flows could happen, such as this:

- Normal -> Problem -> Correction -> Evaluate -> Normal

In this, we found a Problem, made a Correction, Evaluated the correction, and finding things are OK, reset the state back to Normal.

Another example:

- Normal -> Problem -> Correction -> Evaluate -> Escalate -> EscalateWait

Here we run the entire pipeline directly, not solving the issue, alerting on-call, and finally waiting for them.

In the EscalateWait it makes sense to have a "timeout" style action, where if too long has passed (maybe 5-10 minutes), the state pipeline will reset, and go back to Normal.  This allows continued state tracking and potential problem solving.   Maybe the problem resolved itself, or maybe it will work the second time or you get different actions triggering because of the way you set up your action evaluation rules.

## 4. Create Additional State Pipelines to Describe your Bots

Image a web server with the following state pipeline for "Traffic":

Normal -> Low -> High -> None

In a reset state, the Traffic state will be "Normal".  Then actions will run and jump to one of the other states, low, high or none.  Here are a few ways this could happen:

Normal -> None

You had normal traffic, then all of a sudden it stopped, and now you have no traffic.

Normal -> Low -> None

You had normal traffic, then it started to go down, then no traffic.

Normal -> High -> None

You has normal traffic, then it started getting higher, now you have no traffic.

These each have a different pattern to them, don't they?  And different causes may be found because the previous state went from "Low -> None" or "High -> None".

This transition is allowed because "None" comes last.  So dont make a flow like: "Normal -> None -> Low -> High", because the only way to get from High to None is to go back to Normal.  This is a key design feature, and not a flaw.  It will encourage certain thinking about how the states change to get to a resolution, and then re-test from scratch again.

## 5. Add More Queries to your Bot Group

Every Bot in a Bot Group gets all their data from Queries in the Bot Group.  So add all the queries you need.  It's best to use queries you already use in other systems, like Graphana, because you are familiar with them, and can just paste them into Sireus.  It's better to experiment with queries in dedicated query inspection software then once you have them worked out, bring them here as Sireus is not meant to be good at query editing or viewing.

Only bring queries in when you have variables you want to extract from the queries.  You can extract more than 1 variable per query, using the metrics matching, but remember to keep the Bot Key in one of the Metric Names so that you can match the query to all the Bots in the Bot Group.

## 6. Add an State Condition to Cover Every Troubleshooting State

State Conditions are only evalauted if all their Required States are currently active.  So you need to have at least 1 State Condition for every state in your troubleshooting pipeline, which will move the pipeline forward.

You will usually want more than 1 action, so that you can differentiate different kinds of situations inside that state.

## 7. Add an State Condition to Cover Non-Troubleshooting State

For all the non-troubleshooting state pipelines, like the "Traffic" example above, you also want State Conditions to adjust their states.  Always remember you will either need to advance the state, or reset the state.  You can also auto-advance the state, where you don't need to specify the state name, it will just go to the next state.  

Auto-advancing state pipelines is useful for a state like, "Corrections":

- Zero -> One -> Two -> Escalate

Here we can track how many corrections we have performed so far, defaulting to "Zero", and then auto-advance it into the "Corrections.Escalate" state, because we tried 2 things before.

## 8. Make Boolean State Condition Considerations for all your State Conditions

All State Conditions could have at least 1 State Condition Consideration that can be 0 or not-0.  If it's not 0, it will get a score, and as long as that score is above the minimum threshold, it will fire, allowing you to enter different states.

Starting with simple boolean considerations is best.  Use more complex curve based calculations when you are in the deeper tuning phases of working with Sireus.

Boolean variables can be set as 0 (false) and 1 (true) values, and you can test them by multiplying them together:

- service_up * service_has_traffic == 1 * 0 == 0

If the service is up (service_up=1) and there is no traffic (service_has_traffic=0), then 1 * 0 = 0, the consideration score will be 0, and this wont trigger.

In this case, if you want the action to fire, then use the "DecBoolean" instead of "IncBoolean".  "DecBoolean" will be active on 0 value.  "IncBoolean" will be active on 1 value.  So you can flip the values using the Inc (Increasing) or Dec (Decreasing) curves.

## 9. Start with No Commands, Only Exporting Metrics

In the beginning, it is best to just export your states and action scores to your metric system, so you can see how your configuration works over time.  Don't rush to set up a command to run until you have validated how your actions trigger with your current setup.

Use the interactive mode to go back in time and look at different data sets, to see what the actions would be.  Go to previous outages, and check how the actions would peform.  This is the way to check that your changes will give you good results.

If you can match up your scoring with past data accurately, you give yourself the best opportunity to match future data as well.

## 10. Once Confidence in State Condition Score Exists, Add Commands

Start with adding alerting and escalation commands, because these bring a human to check on things.  This is the safest first change.  At worst, you increase your noise, and then can tune your action scores.

Next move to idempotent potential fixes, which should not break anything.

Finally, once you have good confidence add targeted fixes, which can try to fix specific problems.

