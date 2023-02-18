package data

import (
	"encoding/json"
	"errors"
	"time"
)

type (
	// Allows for sorting map[string]float64 by value or key, and also stores a human readable formatted string
	PairFloat64 struct {
		Key       string
		Value     float64
		Formatted string
	}
)

type (
	// Wrapper for PairFloat64 that makes sorting possible
	PairFloat64List []PairFloat64
)

func (p PairFloat64List) Len() int { return len(p) }
func (p PairFloat64List) Less(i, j int) bool {
	// If they aren't the same value, return the lesser
	if p[i].Value != p[j].Value {
		return p[i].Value < p[j].Value
	} else {
		// Else, test their Key names, and return lesser of the strings so order is consistent.  Only for consistency
		return StringCompare(p[i].Key, p[j].Key) > 0
	}
}
func (p PairFloat64List) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type (
	// Allows for sorting BotConditionData by FinalScore
	PairBotConditionData struct {
		Key   string
		Value BotConditionData
	}
)

type PairBotConditionDataList []PairBotConditionData

func (p PairBotConditionDataList) Len() int { return len(p) }
func (p PairBotConditionDataList) Less(i, j int) bool {
	// If they aren't the same value, return the lesser
	if p[i].Value.FinalScore != p[j].Value.FinalScore {
		return p[i].Value.FinalScore < p[j].Value.FinalScore
	} else {
		// Else, test their Key names, and return lesser of the strings so order is consistent.  Only for consistency
		return StringCompare(p[i].Key, p[j].Key) > 0
	}
}
func (p PairBotConditionDataList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type (
	// Cant unmarshal JSON into time.Duration, so wrapping it
	Duration time.Duration
)

// Marshal JSON for our time.Duration wrapper
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

// Unmarshal JSON for our time.Duration wrapper
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
		return nil
	default:
		return errors.New("invalid duration")
	}
}

// Compares 2 strings, -1 is s1 is less, 0 is equal, 1 is s1 is greater
func StringCompare(s1, s2 string) int {
	lens := len(s1)
	if lens > len(s2) {
		lens = len(s2)
	}
	for i := 0; i < lens; i++ {
		if s1[i] != s2[i] {
			return int(s1[i]) - int(s2[i])
		}
	}
	return len(s1) - len(s2)
}
