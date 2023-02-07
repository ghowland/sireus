package appdata

import (
	"sort"
)

type PairBotActionData struct {
	Key   string
	Value BotActionData
}

type PairBotActionDataList []PairBotActionData

func (p PairBotActionDataList) Len() int { return len(p) }
func (p PairBotActionDataList) Less(i, j int) bool {
	return p[i].Value.FinalScore < p[j].Value.FinalScore
}
func (p PairBotActionDataList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func SortMapStringActionDataByFinalScore(input map[string]BotActionData, sortForward bool) PairBotActionDataList {
	pairList := PairBotActionDataList{}

	// Create the PairList
	for key, value := range input {
		newPair := PairBotActionData{
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
