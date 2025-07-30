package datamodel

type OverviewCommand struct {
	Action string `json:"action"` // Unique Key
	Icon   string `json:"icon"`
	Link   string `json:"link"`
	Field  string `json:"field"`
}
