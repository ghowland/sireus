package util

import (
	"sort"
)

type PairFloat64 struct {
	Key       string
	Value     float64
	Formatted string
}

type PairFloat64List []PairFloat64

func (p PairFloat64List) Len() int           { return len(p) }
func (p PairFloat64List) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairFloat64List) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func SortMapStringFloat64ByKey(input map[string]float64) PairFloat64List {
	var keys []string

	for key, _ := range input {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var output PairFloat64List
	for _, key := range keys {
		pair := PairFloat64{
			Key:   key,
			Value: input[key],
		}
		output = append(output, pair)
	}

	return output
}

// TODO(ghowland):PERF: Inefficient, but I am working quickly.  Make it better later.
func SortMapStringFloat64ByValue(input map[string]float64, sortForward bool) PairFloat64List {
	pairList := PairFloat64List{}

	// Create the PairList
	for key, value := range input {
		newPair := PairFloat64{
			Key:   key,
			Value: value,
		}
		pairList = append(pairList, newPair)
	}

	if sortForward {
		sort.Sort(pairList)
	} else {
		sort.Reverse(pairList)
	}

	//log.Printf("Sorted: %v", pairList)

	return pairList
}
