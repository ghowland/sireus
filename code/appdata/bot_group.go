package appdata

type ActionConsideration struct {
	Name       string  `json:"name"`
	Weight     float32 `json:"weight"`
	CurveName  string  `json:"curve"`
	RangeStart float32 `json:"range_start"`
	RangeEnd   float32 `json:"range_end"`
}

type Action struct {
	Name           string                `json:"name"`
	Info           string                `json:"info"`
	Weight         float32               `json:"weight"`
	WeightMin      float32               `json:"weight_min"`
	Considerations []ActionConsideration `json:"considerations"`
}

type Bot struct {
	Name    string   `json:"name"`
	Info    string   `json:"info"`
	Actions []Action `json:"actions"`
}
