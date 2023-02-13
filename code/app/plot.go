package app

import (
	"encoding/json"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/util"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

// Returns JSON data needed to create a Plotly graph for our Curves
func GetAPIPlotData(c *fiber.Ctx) string {
	input := util.ParseContextBody(c)
	//log.Println("Get API Plot Data: ", input)

	if input["name"] == "" {
		failureResult := map[string]interface{}{
			"_failure": "Name not found, aborting",
		}
		failureJson, _ := json.Marshal(failureResult)
		return string(failureJson)
	}

	curveData := LoadCurveData(data.SireusData.AppConfig, input["name"])

	mapData := map[string]interface{}{
		"title":  curveData.Name,
		"plot_x": GetCurveDataX(curveData),
		"plot_y": curveData.Values,
	}

	xPos, err := strconv.ParseFloat(input["x"], 32)
	util.Check(err)

	if xPos >= 0 {
		mapData["plot_selected_x"] = xPos
		mapData["plot_selected_y"] = GetCurveValue(curveData, float64(xPos))
	}

	jsonOutput, _ := json.Marshal(mapData)
	jsonString := string(jsonOutput)

	//log.Println("Get API Plot Result: ", json_string)

	return jsonString
}

func GetRawMetricsJSON(c *fiber.Ctx) string {
	queryKey := c.Query("query_key")

	queryResult, ok := GetQueryResultByQueryKey(&data.SireusData.Site, queryKey)
	if !ok {
		return "{}"
	}

	return util.PrintJson(queryResult)
}

func GetAPIPlotMetrics(c *fiber.Ctx) string {
	input := util.ParseContextBody(c)
	log.Println("Get API Plot Metrics: ", util.PrintJson(input))

	queryKey := input["query_key"]

	queryResult, ok := GetQueryResultByQueryKey(&data.SireusData.Site, queryKey)
	if !ok {
		return "{}"
	}

	xArray := []float64{0, 0}
	yArray := []float64{0, 1}

	//queryResult.Result.PrometheusResponse.Data.Result[0].Values

	mapData := map[string]interface{}{
		"title":  queryResult.Query,
		"plot_x": xArray,
		"plot_y": yArray,
	}

	//xPos, err := strconv.ParseFloat(input["x"], 32)
	//util.Check(err)
	//
	//if xPos >= 0 {
	//	mapData["plot_selected_x"] = xPos
	//	mapData["plot_selected_y"] = GetCurveValue(curveData, float64(xPos))
	//}

	jsonOutput, _ := json.Marshal(mapData)
	jsonString := string(jsonOutput)

	//log.Println("Get API Plot Result: ", json_string)

	return jsonString
}
