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

// Query the Prometheus metric server
func QueryPrometheus(host string, port int, queryType data.BotQueryType, query string, timeStart time.Time, duration time.Duration) data.PrometheusResponse {
	queryStartTime := util.GetTimeNow()

	start := timeStart.UTC().Format(time.RFC3339)

	durationSeconds := duration.Seconds()

	end := timeStart.UTC().Add(time.Second * time.Duration(durationSeconds)).Format(time.RFC3339)

	requestUrl := fmt.Sprintf("http://%s:%d/api/v1/%s?query=%s&start=%s&end=%s&step=15s", host, port, queryType.String(), url.QueryEscape(query), start, end)

	var jsonResponse data.PrometheusResponse

	// Set the time, so we know when we got it
	jsonResponse.RequestURL = requestUrl
	jsonResponse.RequestTime = queryStartTime
	jsonResponse.ResponseTime = util.GetTimeNow()

	resp, err := http.Get(requestUrl)
	if util.Check(err) {
		jsonResponse.IsError = true
		jsonResponse.ErrorMessage = fmt.Sprintf("Couldn't fetch Prometheus URL: %v   URL: %s", err, requestUrl)
		log.Printf(jsonResponse.ErrorMessage)
		return jsonResponse
	}

	body, err := io.ReadAll(resp.Body)
	util.CheckLog(err)
	if util.Check(err) {
		jsonResponse.IsError = true
		jsonResponse.ErrorMessage = fmt.Sprintf("Couldn't read Prometheus body: %v   URL: %s", err, requestUrl)
		log.Printf(jsonResponse.ErrorMessage)
		return jsonResponse
	}

	err = json.Unmarshal(body, &jsonResponse)
	util.CheckLog(err)
	if util.Check(err) {
		jsonResponse.IsError = true
		jsonResponse.ErrorMessage = fmt.Sprintf("Couldn't unmarshall JSON: %v   URL: %s", err, requestUrl)
		log.Printf(jsonResponse.ErrorMessage)
		return jsonResponse
	}

	return jsonResponse
}

// Extract our ephemeral Bots from the Prometheus response, using the BotKey extractor information
func ExtractBotsFromPromData(response data.PrometheusResponse, botGroup *data.BotGroup) []data.Bot {
	bots := make(map[string]data.Bot)

	for _, resultItem := range response.Data.Result {
		name := resultItem.Metric[botGroup.BotExtractor.Key]

		_, exists := bots[name]
		if !exists {
			bots[name] = data.Bot{
				Name:           name,
				LockKey:        fmt.Sprintf("%s.%s", botGroup.LockKey, name),
				ActionData:     map[string]data.BotActionData{},
				StateValues:    []string{},
				VariableValues: map[string]float64{},
			}
		}
	}

	// Add all the bots to a final array.  The map allowed us to ensure no duplicate entries, as that is allowed.
	var botArray []data.Bot
	for _, bot := range bots {
		botArray = append(botArray, bot)
	}

	//log.Print("Bots: ", botArray)
	return botArray
}
