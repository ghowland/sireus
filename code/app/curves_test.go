package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// NOTE(ghowland): Have to go back 2 directories here, and also in the config data, so need a separate config
var testAppConfigPath = "../../config/test_config.json"

func TestSomething(t *testing.T) {
	appConfig := LoadConfig(testAppConfigPath)

	assert.NotNil(t, appConfig, "App Config found, loaded, not nil")

	curve_data, _ := LoadCurveData("inc_smooth")

	assert.NotNil(t, curve_data, "Curve data found, loaded, not nil")
	assert.Equal(t, curve_data.Values[0], float64(0), "First value of this curve should be 0")

}
