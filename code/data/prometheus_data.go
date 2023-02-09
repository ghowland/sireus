package data

import "time"

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
