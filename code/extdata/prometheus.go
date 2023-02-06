package extdata

import (
	"encoding/json"
	"fmt"
	"github.com/ghowland/sireus/code/appdata"
	"github.com/ghowland/sireus/code/util"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type PrometheusResponseDataResult struct {
	Metric map[string]string `json:"metric"`
	Values [][]interface{}   `json:"values"`
}

type PrometheusResponseData struct {
	ResultType string                         `json:"resultType"`
	Result     []PrometheusResponseDataResult `json:"result"`
}

type PrometheusResponse struct {
	Status string                 `json:"status"`
	Data   PrometheusResponseData `json:"data"`
}

func QueryPrometheus(host string, port int, queryType appdata.BotQueryType, query string, timeStart time.Time, duration int) PrometheusResponse {
	start := timeStart.UTC().Format(time.RFC3339)

	end := timeStart.UTC().Add(time.Second * time.Duration(duration)).Format(time.RFC3339)

	url := fmt.Sprintf("http://%s:%d/api/v1/%s?query=%s&start=%s&end=%s&step=15s", host, port, queryType.String(), url.QueryEscape(query), start, end)

	log.Print("Prom URL: ", url)

	resp, err := http.Get(url)
	util.Check(err)

	body, err := io.ReadAll(resp.Body)
	util.Check(err)

	//log.Print("Prom Result: ", string(body))

	var jsonResponse PrometheusResponse
	err = json.Unmarshal(body, &jsonResponse)
	util.Check(err)

	return jsonResponse
}

func ExtractBotsFromPromData(data PrometheusResponse, botKey string) []appdata.Bot {
	bots := make(map[string]appdata.Bot)

	for _, resultItem := range data.Data.Result {
		name := resultItem.Metric[botKey]

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

	//log.Print("Bots: ", botArray)

	return botArray
}
