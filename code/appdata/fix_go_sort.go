package appdata

import (
	"github.com/ghowland/sireus/code/data"
	"sort"
)

func SortMapStringActionDataByFinalScore(input map[string]data.BotActionData, sortForward bool) data.PairBotActionDataList {
	pairList := data.PairBotActionDataList{}

	// Create the PairList
	for key, value := range input {
		newPair := data.PairBotActionData{
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
