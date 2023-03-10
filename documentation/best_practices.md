# Best Practices

- [Best Practices](#best-practices)
  * [Naming your State Conditions](#naming-your-state-conditions)
  * [Setting the State Condition Consideration Weights](#setting-the-state-condition-consideration-weights)

### Naming your State Conditions

TL;DR - Name your State Conditions to answer "What condition is being uniquely described here?" using CSV of descriptions of the evaluation elements.

I would recommend naming your State Conditions to describe the state they represent.  This name would answer the question, "What is going on that caused this condition to be true?"

**Example names:**

 - Service Stopped, Out of Space
 - Service Stopped, Has Storage

These names answer the question: **Why perform this State Condition?**

With "*Service Stopped, Out of Space*", it is likely being out of storage is what caused the service to stop, so an State Condition will be executed to try to deal with that.

With "*Service Stopped, Has Storage*", we know the service is stopped, but it still has storage, so we want to run a different command that deals with problems unrelated to running out of storage.

This is a simple 2 set problem, but let's expand the list to see why this is a scalable naming pattern:

 - Service Stopped, Out of Space
 - Service Stopped, Has Storage
 - Service Stopped, Won't Restart
 - Service Running, Many Errors
 - Service Running, Many Timeouts
 - Service Running, Database Connection Errors
 - Service Running, Too Busy
 - Service Running, Maybe Under Attack

At this point if we were naming things differently, it would become hard to add more State Conditions and understand what the differences of them are.  This in some way just turns the evaluations into text, but it should also simplify those evaluations into big picture concepts so that new users can get a grasp on things easier, and experienced used can quickly differentiate.

This still has issues in that you can have more than 2 conditions.  For this, consider using Synthethic Variables to create a combination of values so that you can test them as a boolean.  In this way as you grow in variables, you can reduce them into Synthethic Variables to keep the State Condition evaluation logic simpler, and the names easier to read and understand, even as the number of State Conditions continues to increase.

### Setting the State Condition Consideration Weights

TL;DR - Keep State Condition Consideration weights between 0.1 and 10.0, with most being closer to 1.0.

Consideration weights should try to stay in the low numbers, the best weight being \~1.0.  Anything under 10 would be good for a particular strong weight to offset the more normal 1.0 weights, as a lower importance weight could be given 0.5 or 0.2 as it's Consideration Weight.

The reason for this is that in the calculations, there is a running score that multiplies all the scores together.  If a lot of \~1.0 values are multiplying each other, then the final result will be in the \~1.0 range.  If there are many different values such as "500, 10, 30, 1.0, 2000", these numbers are so different your Consideration Final Score will be very difficult to understand or control.

Using numbers such as "1.0, 1.3, 0.7, 1.0, 2.5" allows a set of Considerations that have relative importance to each other without swingingly wildly out of control if a one or more value is at an extreme.

The State Condition Weight is where you differentiate State Conditions from each other by their weighted State Condition Final Score, which is the multiplier of the result of the Consideration Final Score by the State Condition Weight.  This allows you to increase an State Condition by 500x, 2000x, etc.  If you want to put State Conditions into different categories of priority this way, do it with State Condition Weight, and leave the Consideration Weight near 1.0.