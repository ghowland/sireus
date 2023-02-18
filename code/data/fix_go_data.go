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
		return p[i].Key < p[j].Key
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
	return p[i].Value.FinalScore < p[j].Value.FinalScore
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
