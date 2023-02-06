package appdata

import (
	"encoding/json"
	"github.com/ghowland/sireus/code/util"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetAPIPlotData(appConfig AppConfig, c *fiber.Ctx) string {
	input := util.ParseContextBody(c)
	//log.Println("Get API Plot Data: ", input)

	if input["name"] == "" {
		failureResult := map[string]interface{}{
			"_failure": "Name not found, aborting",
		}
		failureJson, _ := json.Marshal(failureResult)
		return string(failureJson)
	}

	curveData := LoadCurveData(appConfig, input["name"])

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
