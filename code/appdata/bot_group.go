package appdata

import (
	"encoding/json"
	"fmt"
	"github.com/ghowland/sireus/code/util"
	"os"
)

type ActionConsideration struct {
	Name       string  `json:"name"`
	Weight     float32 `json:"weight"`
	CurveName  string  `json:"curve"`
	RangeStart float32 `json:"range_start"`
	RangeEnd   float32 `json:"range_end"`
}

type Action struct {
	Name           string                `json:"name"`
	Info           string                `json:"info"`
	Weight         float32               `json:"weight"`
	WeightMin      float32               `json:"weight_min"`
	Considerations []ActionConsideration `json:"considerations"`
}

type Bot struct {
	Name    string   `json:"name"`
	Info    string   `json:"info"`
	Actions []Action `json:"actions"`
}

type CurveData struct {
	Name   string    `json:"name"`
	Values []float32 `json:"values"`
}

func LoadCurveData(app_config AppConfig, name string) CurveData {
	path := fmt.Sprintf(app_config.CurvePathFormat, name)

	curveData, err := os.ReadFile((path))
	util.Check(err)

	var curve_data CurveData
	json.Unmarshal([]byte(curveData), &curve_data)

	//log.Println("Load Curve Data: ", curve_data.Name)

	return curve_data
}
