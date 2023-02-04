package appdata

import (
	"encoding/json"
	"github.com/ghowland/sireus/code/util"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetAPIPlotData(app_config AppConfig, c *fiber.Ctx) string {
	input := util.ParseContextBody(c)
	//log.Println("Get API Plot Data: ", input)

	if input["name"] == "" {
		failure_result := map[string]interface{}{
			"_failure": "Name not found, aborting",
		}
		failure_json, _ := json.Marshal(failure_result)
		return string(failure_json)
	}

	curve_data := LoadCurveData(app_config, input["name"])

	map_data := map[string]interface{}{
		"title":  curve_data.Name,
		"plot_x": GetCurveDataX(curve_data),
		"plot_y": curve_data.Values,
	}

	x_pos, err := strconv.ParseFloat(input["x"], 32)
	util.Check(err)

	if x_pos >= 0 {
		map_data["plot_selected_x"] = x_pos
		map_data["plot_selected_y"] = GetCurveValue(curve_data, float32(x_pos))
	}

	json_output, _ := json.Marshal(map_data)
	json_string := string(json_output)

	//log.Println("Get API Plot Result: ", json_string)

	return json_string
}
