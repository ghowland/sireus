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

func QueryPrometheus(host string, port int, queryType appdata.BotQueryType, query string, timeStart time.Time, duration int) map[string]interface{} {
	start := timeStart.UTC().Format(time.RFC3339)

	end := timeStart.UTC().Add(time.Second * time.Duration(duration)).Format(time.RFC3339)

	url := fmt.Sprintf("http://%s:%d/api/v1/%s?query=%s&start=%s&end=%s&step=15s", host, port, queryType.String(), query, start, end)

	log.Print("Prom URL: ", url)

	resp, err := http.Get(url)
	util.Check(err)

	body, err := io.ReadAll(resp.Body)
	util.Check(err)

	//log.Print("Prom Result: ", string(body))

	var jsonResultInt interface{}
	err = json.Unmarshal(body, &jsonResultInt)
	util.Check(err)
	jsonResult := jsonResultInt.(map[string]interface{})

	return jsonResult
}

func ExtractBotsFromPromData(data map[string]interface{}, botKey string) []appdata.Bot {
	//log.Print("Extra From: ", data)

	bots := make(map[string]appdata.Bot)

	resultItems := data["data"].(map[string]interface{})["result"].([]interface{})

	for _, resultItem := range resultItems {
		item := resultItem.(map[string]interface{})
		metric := item["metric"].(map[string]interface{})
		//log.Print("Item: ", metric)

		name := metric[botKey].(string)

		_, exists := bots[name]
		if !exists {
			bots[name] = appdata.Bot{
				Name: name,
			}
		}
	}
	//log.Print("Bots: ", bots)

	// Add all the bots to a final array.  The map allowed us to ensure no duplicate entries, as that is allowed.
	botArray := []appdata.Bot{}
	for _, bot := range bots {
		botArray = append(botArray, bot)
	}

	log.Print("Bots: ", botArray)

	return botArray
}
