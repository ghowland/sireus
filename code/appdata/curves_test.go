package appdata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// NOTE(ghowland): Have to go back 2 directories here, and also in the config data, so need a separate config
var test_app_config_path = "../../config/test_config.json"

func TestSomething(t *testing.T) {
	app_config := LoadConfig(test_app_config_path)

	assert.NotNil(t, app_config, "App Config found, loaded, not nil")

	curve_data := LoadCurveData(app_config, "inc_smooth")

	assert.NotNil(t, curve_data, "Curve data found, loaded, not nil")
	assert.Equal(t, curve_data.Values[0], float32(0), "First value of this curve should be 0")

}
