package appdata

import (
	"encoding/json"
	"fmt"
	"github.com/ghowland/sireus/code/util"
	"os"
)

type CurveData struct {
	Name   string    `json:"name"`
	Values []float32 `json:"values"`
}

func LoadCurveData(appConfig AppConfig, name string) CurveData {
	path := fmt.Sprintf(appConfig.CurvePathFormat, name)

	curveDataInput, err := os.ReadFile(path)
	util.Check(err)

	var curveData CurveData
	err = json.Unmarshal(curveDataInput, &curveData)
	util.Check(err)

	//log.Println("Load Curve Data: ", curve_data.Name)

	return curveData
}

func GetCurveDataX(curveData CurveData) []float32 {
	var xArray []float32

	for i := 0; i < len(curveData.Values); i++ {
		xArray = append(xArray, float32(i)*0.01)
	}

	return xArray
}

func GetCurveValue(curveData CurveData, x float32) float32 {

	for i := 0; i < len(curveData.Values); i++ {
		curPosX := float32(i) * 0.01
		if x <= curPosX {
			return curveData.Values[i]
		}
	}

	return curveData.Values[len(curveData.Values)-1]
}
