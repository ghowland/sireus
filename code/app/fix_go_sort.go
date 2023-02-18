package app

import (
	"github.com/ghowland/sireus/code/data"
	"sort"
)

// Go doesnt handle map sorting easily, so this is the fix-up
func SortMapStringConditionDataByFinalScore(input map[string]data.BotConditionData, sortForward bool) data.PairBotConditionDataList {
	pairList := data.PairBotConditionDataList{}

	// Create the PairList
	for key, value := range input {
		newPair := data.PairBotConditionData{
			Key:   key,
			Value: value,
		}
		pairList = append(pairList, newPair)
	}

	if sortForward {
		sort.Sort(pairList)
	} else {
		sort.Sort(sort.Reverse(pairList))
	}

	//log.Printf("Sorted: %v", pairList)

	return pairList
}
