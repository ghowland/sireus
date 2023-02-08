package app

import (
	"encoding/json"
	"fmt"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/util"
	"os"
)

type CurveData struct {
	Name   string    `json:"name"`
	Values []float64 `json:"values"`
}

func LoadCurveData(appConfig data.AppConfig, name string) CurveData {
	path := fmt.Sprintf(appConfig.CurvePathFormat, name)

	curveDataInput, err := os.ReadFile(path)
	util.Check(err)

	var curveData CurveData
	err = json.Unmarshal(curveDataInput, &curveData)
	util.Check(err)

	//log.Println("Load Curve Data: ", curve_data.Name)

	return curveData
}

func GetCurveDataX(curveData CurveData) []float64 {
	var xArray []float64

	for i := 0; i < len(curveData.Values); i++ {
		xArray = append(xArray, float64(i)*0.01)
	}

	return xArray
}

func GetCurveValue(curveData CurveData, x float64) float64 {

	for i := 0; i < len(curveData.Values); i++ {
		curPosX := float64(i) * 0.01
		if x <= curPosX {
			return curveData.Values[i]
		}
	}

	return curveData.Values[len(curveData.Values)-1]
}
