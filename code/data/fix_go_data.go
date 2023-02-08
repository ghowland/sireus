package data

import (
	"encoding/json"
	"errors"
	"time"
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

type (
	// Cant unmarshal JSON into time.Duration, so wrapping it
	Duration time.Duration
)

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

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
