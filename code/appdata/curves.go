package appdata

func GetCurveDataX(curve_data CurveData) []float32 {
	var x_array []float32

	for i := 0; i < len(curve_data.Values); i++ {
		x_array = append(x_array, float32(i)*0.01)
	}

	return x_array
}

func GetCurveValue(curve_data CurveData, x float32) float32 {

	for i := 0; i < len(curve_data.Values); i++ {
		cur_pos_x := float32(i) * 0.01
		if x <= cur_pos_x {
			return curve_data.Values[i]
		}
	}

	return curve_data.Values[len(curve_data.Values)-1]
}
