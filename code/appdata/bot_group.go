package appdata

type ActionConsideration struct {
	Name       string  `json:"name"`
	Weight     float32 `json:"weight"`
	CurveName  string  `json:"curve"`
	RangeStart float32 `json:"range_start"`
	RangeEnd   float32 `json:"range_end"`
}

type ActionCommand struct {
}

type Action struct {
	Name           string                `json:"name"`
	Info           string                `json:"info"`
	Weight         float32               `json:"weight"`
	WeightMin      float32               `json:"weight_min"`
	Considerations []ActionConsideration `json:"considerations"`
	Command        ActionCommand
}

type Bot struct {
	Name    string   `json:"name"`
	Info    string   `json:"info"`
	Actions []Action `json:"actions"`
}

type BotQueryType int64

const (
	Undefined BotQueryType = iota
	PrometheusQueryRange
)

func (bqt BotQueryType) String() string {
	switch bqt {
	case Undefined:
		return "Undefined"
	case PrometheusQueryRange:
		return "PromQueryRange"
	}
	return "Unknown"
}

type BotQuery struct {
	Type BotQueryType `json:"query_type"`
	Name string       `json:"name"`
	Info string       `json:"info"`
}

type BotGroup struct {
	Name    string   `json:"name"`
	Info    string   `json:"info"`
	Actions []Action `json:"actions"`
	Bots    []Bot
}

type Site struct {
	Name      string   `json:"name"`
	Paths     []string `json:"paths"`
	BotGroups []BotGroup
}
