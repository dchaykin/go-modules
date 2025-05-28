package datamodel

type AccessConfig struct {
	Partner   string  `json:"partner"`
	Algorithm *string `json:"algo,omitempty"`
}

type OverviewCommand struct {
	Action string `json:"action"` // Unique Key
	Icon   string `json:"icon"`
	Link   string `json:"link"`
	Field  string `json:"field"`
}
