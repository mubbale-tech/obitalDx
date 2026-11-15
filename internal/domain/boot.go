package domain

import "encoding/json"

type InitialBootPayload struct {
	Evaluation EvaluationResult `json:"evaluation"`
	Missions   []MissionProfile `json:"missions"`
	Hazards    []NASAHazard     `json:"hazards"`
	Scenarios  []ScenarioPreset `json:"scenarios"`
}

func MustJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(b)
}
