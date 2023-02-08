package extdata

import (
	"encoding/json"
	"fmt"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/util"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type (
	// Data inside the payload of the PrometheusResponseData
	PrometheusResponseDataResult struct {
		Metric map[string]string `json:"metric"`
		Values [][]interface{}   `json:"values"`
	}
)

type (
	// Payload for PrometheusResponse
	PrometheusResponseData struct {
		ResultType string                         `json:"resultType"`
		Result     []PrometheusResponseDataResult `json:"result"`
	}
)

type (
	// Response from Prometheus.  I made a short-hand version of this instead of using the one from Prometheus for convenience.
	PrometheusResponse struct {
		Status       string                 `json:"status"`
		Data         PrometheusResponseData `json:"data"`
		RequestTime  time.Time              // When the Request was made
		ResponseTime time.Time              // When the Response was received
	}
)

type (
	// A single Query result
	QueryResult struct {
		QueryServer        string // Server this Query came from
		QueryType          data.BotQueryType
		QueryName          string             // The Query
		PrometheusResponse PrometheusResponse // The Response
	}
)

type (
	// Stores all our QueryResults
	QueryManager struct {
		Results []QueryResult
	}
)

// Query the Prometheus metric server
func QueryPrometheus(host string, port int, queryType data.BotQueryType, query string, timeStart time.Time, duration int) PrometheusResponse {
	queryStartTime := time.Now()

	start := timeStart.UTC().Format(time.RFC3339)

	end := timeStart.UTC().Add(time.Second * time.Duration(duration)).Format(time.RFC3339)

	requestUrl := fmt.Sprintf("http://%s:%d/api/v1/%s?query=%s&start=%s&end=%s&step=15s", host, port, queryType.String(), url.QueryEscape(query), start, end)

	log.Print("Prom URL: ", requestUrl)

	resp, err := http.Get(requestUrl)
	util.Check(err)

	body, err := io.ReadAll(resp.Body)
	util.Check(err)

	//log.Print("Prom Result: ", string(body))

	var jsonResponse PrometheusResponse
	err = json.Unmarshal(body, &jsonResponse)
	util.Check(err)

	// Set the time, so we know when we got it
	jsonResponse.RequestTime = queryStartTime
	jsonResponse.ResponseTime = time.Now()

	return jsonResponse
}

// Extract our ephemeral Bots from the Prometheus response, using the BotKey extractor information
func ExtractBotsFromPromData(response PrometheusResponse, botKey string) []data.Bot {
	bots := make(map[string]data.Bot)

	for _, resultItem := range response.Data.Result {
		name := resultItem.Metric[botKey]

		_, exists := bots[name]
		if !exists {
			bots[name] = data.Bot{
				Name:        name,
				ActionData:  map[string]data.BotActionData{},
				StateValues: []string{},
			}
		}
	}
	//log.Print("Bots: ", bots)

	// Add all the bots to a final array.  The map allowed us to ensure no duplicate entries, as that is allowed.
	var botArray []data.Bot
	for _, bot := range bots {
		botArray = append(botArray, bot)
	}

	//log.Print("Bots: ", botArray)

	return botArray
}
