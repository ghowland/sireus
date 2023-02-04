package extdata

import (
	"encoding/json"
	"fmt"
	"github.com/ghowland/sireus/code/appdata"
	"github.com/ghowland/sireus/code/util"
	"io"
	"log"
	"net/http"
	"time"
)

func QueryPrometheus(host string, port int, query string, time_start time.Time, duration int) map[string]interface{} {
	start := time_start.UTC().Format(time.RFC3339)

	end := time_start.UTC().Add(time.Second * time.Duration(duration)).Format(time.RFC3339)

	//start := time_start.Format()

	url := fmt.Sprintf("http://%s:%d/api/v1/%s&start=%s&end=%s&step=15s", host, port, query, start, end)

	log.Print("Prom URL: ", url)

	resp, err := http.Get(url)
	util.Check(err)

	body, err := io.ReadAll(resp.Body)
	util.Check(err)

	//log.Print("Prom Result: ", string(body))

	var json_result_int interface{}
	err = json.Unmarshal(body, &json_result_int)
	util.Check(err)
	json_result := json_result_int.(map[string]interface{})

	return json_result
}

func ExtractBotsFromPromData(data map[string]interface{}, bot_key string) map[string]appdata.Bot {
	//log.Print("Extra From: ", data)

	bots := make(map[string]appdata.Bot)

	result_items := data["data"].(map[string]interface{})["result"].([]interface{})

	for _, result_item := range result_items {
		item := result_item.(map[string]interface{})
		metric := item["metric"].(map[string]interface{})
		//log.Print("Item: ", metric)

		name := metric[bot_key].(string)

		_, exists := bots[name]
		if !exists {
			bots[name] = appdata.Bot{
				Name: name,
			}
		}

		//log.Print("Bot: ", name)
	}

	log.Print("Bots: ", bots)

	return bots
}
