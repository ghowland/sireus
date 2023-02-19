package app

import (
	"errors"
	"fmt"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strings"
)

// For dynamic metric keys, this will clean them for Prometheus
func CleanMetricKeyString(key string) string {
	key = strings.ToLower(key)
	key = util.StringReplaceUnsafeChars(key, " [](){}=-!@#$%^&*()+<>,./?;:'\"`~", "_")
	for strings.Contains(key, "__") {
		key = strings.Replace(key, "__", "_", -1)
	}
	return key
}

func IsEqualMapStringString(m1 map[string]string, m2 map[string]string) bool {
	// Must be the same size maps to be equal
	if len(m1) != len(m2) {
		return false
	}

	for key1, value1 := range m1 {
		value2, ok := m2[key1]
		if !ok || value1 != value2 {
			return false
		}
	}

	return true
}

// Get an existing Metric Counter
func GetMetricGauge(key string, labels map[string]string) (*data.PrometheusMetricGauge, error) {
	for index := range data.SireusData.MetricExport.Gauges {
		metric := &data.SireusData.MetricExport.Gauges[index]
		if metric.Key == key && IsEqualMapStringString(metric.Labels, labels) {
			return metric, nil
		}
	}
	return &data.PrometheusMetricGauge{}, errors.New(fmt.Sprintf("Missing Gauge: %s", key))
}

// Get an existing Metric Counter
func GetMetricCounter(key string, labels map[string]string) (*data.PrometheusMetricCounter, error) {
	for index := range data.SireusData.MetricExport.Counters {
		metric := &data.SireusData.MetricExport.Counters[index]
		if metric.Key == key && IsEqualMapStringString(metric.Labels, labels) {
			return metric, nil
		}
	}
	return &data.PrometheusMetricCounter{}, errors.New(fmt.Sprintf("Missing Counter: %s", key))
}

// Set a Metric Gauge, and create it if it doesn/t already exist
func SetMetricGauge(key string, value float64, info string, labels map[string]string) {
	gauge, err := GetMetricGauge(key, labels)
	if util.Check(err) {
		gauge = &data.PrometheusMetricGauge{
			Key:    key,
			Labels: labels,
			Metric: promauto.NewGauge(prometheus.GaugeOpts{
				Name:        key,
				Help:        info,
				ConstLabels: labels,
			}),
		}
		gauge.Metric.Set(value)
		data.SireusData.MetricExport.Gauges = append(data.SireusData.MetricExport.Gauges, *gauge)
	} else {
		gauge.Metric.Set(value)
	}
}

// Set a Metric Counter, and create it if it doesn/t already exist
func AddToMetricCounter(key string, value float64, info string, labels map[string]string) {
	counter, err := GetMetricCounter(key, labels)
	if util.Check(err) {
		counter = &data.PrometheusMetricCounter{
			Key:    key,
			Labels: labels,
			Metric: promauto.NewCounter(prometheus.CounterOpts{
				Name:        key,
				Help:        info,
				ConstLabels: labels,
			}),
		}
		counter.Metric.Add(value)
		data.SireusData.MetricExport.Counters = append(data.SireusData.MetricExport.Counters, *counter)
	} else {
		counter.Metric.Add(value)
	}
}

// Returns the map used for Labels in a Metric, for a Bot
func GetMetricLabelsAndInfo_Bot(botGroup *data.BotGroup, bot *data.Bot) map[string]string {
	labels := map[string]string{
		"service":   "sireus",
		"bot":       bot.Name,
		"bot_group": botGroup.Name,
	}
	return labels
}

// Returns the map used for Labels in a Metric, for a Condition
func GetMetricLabelsAndInfo_Condition(botGroup *data.BotGroup, bot *data.Bot, condition data.Condition) map[string]string {
	labels := map[string]string{
		"service":   "sireus",
		"bot":       bot.Name,
		"bot_group": botGroup.Name,
		"condition": condition.Name,
		"states":    util.PrintStringArrayCSV(condition.RequiredStates),
	}
	return labels
}

// Returns the map used for Labels in a Metric, for a Bot's Variable
func GetMetricLabelsAndInfo_BotVariable(botGroup *data.BotGroup, bot *data.Bot, varName string) map[string]string {
	labels := map[string]string{
		"service":   "sireus",
		"bot":       bot.Name,
		"bot_group": botGroup.Name,
		"variable":  varName,
	}
	return labels
}
