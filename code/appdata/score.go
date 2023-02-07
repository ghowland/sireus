package appdata

import (
	"github.com/ghowland/sireus/code/util"
)

func CalculateScore(action Action, actionData BotActionData) float64 {
	var runningScore float64 = 1
	var considerCount int = 0

	for considerName, considerScore := range actionData.ConsiderationFinalScores {
		// We will use a "modified average" to create a calculated score for all the considerations, so need a count
		considerCount++

		//TODO: Process the curve

		// Any Consideration that is 0, means the entire Score is 0, and will never be executed.  It is not Invalid
		if considerScore == 0 {
			return 0
		}

		consider, err := GetActionConsideration(action, considerName)
		if util.Check(err) {
			return 0
		}

		// Multiply the raw Score by the Weight
		rangedScore := util.RangeMapper(considerScore, consider.RangeStart, consider.RangeEnd)

		weightedScore := rangedScore * consider.Weight

		// Move a constantly Running Score
		runningScore *= weightedScore

		//log.Printf("Consider: %s  Score: %.2f  Ranged Score: %.2f  Weighted: %.2f  Running: %.2f", consider.Name, considerScore, rangedScore, weightedScore, runningScore)
	}

	// Mix the numbers together in a "modified average" which yields a good result, especially for low or 0-1 numbers
	calculatedScore := AverageAndFixup(runningScore, considerCount)

	//log.Printf("Calculate: %s  Consider Count: %d  Calc Score: %.2f", action.Name, considerCount, calculatedScore)

	return calculatedScore
}

func AverageAndFixup(runningScore float64, considerCount int) float64 {
	// No considerations is always 0.  We will be dividing by considerCount later...
	if considerCount == 0 {
		return 0
	}

	// Create the modification factor for our averaging function
	var modFactor float64 = 1.0 - (1.0 / float64(considerCount))

	// This is our fudge, that makes the numbers look better.  Especially if Consideration Scores aim between 0-1.
	//NOTE(g): Use Action.Weight to make them much higher for final sorting.  Here in Consideration-land they are better
	//		   between 0 and 1.  It makes reasoning about them easier as well, as all the ranges are normalized.
	//		   Highly recommended to keep the Consideration scores at 0-1, or low numbers like 10 at the highest, then
	//		   use the Action.Weight to massively modify the numbers into different ranges for the Final Score sort.
	var makeUpValue float64 = (1.0 - runningScore) * modFactor

	// Apply the average and fixup to the running score
	var finalScore float64 = runningScore + (makeUpValue * runningScore)

	//log.Printf("AvgFixup: Running Score: %.2f  Count: %d  Mod: %0.2f  Make Up: %.2f  Final Score: %.2f", runningScore, considerCount, modFactor, makeUpValue, finalScore)

	return finalScore
}
