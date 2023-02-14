package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/util"
	"os"
)

type (
	// Points to create a curve.  Standard is 0-1 at 0.1 steps, so 1000 points
	CurveData struct {
		Name   string    `json:"name"`
		Values []float64 `json:"values"`
	}
)

var (
	Curves []CurveData
)

func GetCurve(name string) (CurveData, error) {
	for _, curve := range Curves {
		if curve.Name == name {
			return curve, nil
		}
	}

	curve, err := LoadCurveData(name)
	if util.Check(err) {
		return CurveData{}, errors.New(fmt.Sprintf("Could not find curve: %s", name))
	}

	return curve, nil
}

// Load the Curve data off the disk
func LoadCurveData(name string) (CurveData, error) {
	path := fmt.Sprintf(data.SireusData.AppConfig.CurvePathFormat, name)

	curveDataInput, err := os.ReadFile(path)
	util.Check(err)

	var curveData CurveData
	err = json.Unmarshal(curveDataInput, &curveData)
	if util.CheckNoLog(err) {
		return CurveData{}, errors.New(fmt.Sprintf("Couldnt find curve: %s", name))
	}

	//log.Println("Load Curve Data: ", curve_data.Name)

	return curveData, nil
}

// Get all X axis values, which is just the step from 0-1 at 0.1 intervals
func GetCurveDataX(curveData CurveData) []float64 {
	var xArray []float64

	for i := 0; i < len(curveData.Values); i++ {
		xArray = append(xArray, float64(i)*0.01)
	}

	return xArray
}

// Get the Y value, at an X position, in the curve
func GetCurveValue(curveData CurveData, x float64) float64 {

	for i := 0; i < len(curveData.Values); i++ {
		curPosX := float64(i) * 0.01
		if x <= curPosX {
			return curveData.Values[i]
		}
	}

	return curveData.Values[len(curveData.Values)-1]
}
