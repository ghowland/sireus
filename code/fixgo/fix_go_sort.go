package fixgo

import (
	"github.com/ghowland/sireus/code/data"
	"sort"
)

// Sort a map[string]float64 by their key.  Returns a custom pair list
func SortMapStringFloat64ByKey(input map[string]float64) data.PairFloat64List {
	var keys []string

	for key := range input {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var output data.PairFloat64List
	for _, key := range keys {
		pair := data.PairFloat64{
			Key:   key,
			Value: input[key],
		}
		output = append(output, pair)
	}

	return output
}

// Sort a map[string]float64 by their value.  Returns a custom pair list
// TODO(ghowland):PERF: Inefficient, but I am working quickly.  Make it better later.
func SortMapStringFloat64ByValue(input map[string]float64, sortForward bool) data.PairFloat64List {
	pairList := data.PairFloat64List{}

	// Create the PairList
	for key, value := range input {
		newPair := data.PairFloat64{
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
