package overview

import "github.com/dchaykin/go-modules/datamodel"

const (
	ActionCreate   = "create"
	ActionOpen     = "open"
	ActionPrintPdf = "printPdf"
	ActionDelete   = "delete"
)

type OverviewAction struct {
	Command string `json:"command"`
	Icon    string `json:"icon"`
	Link    string `json:"link"`
	Field   string `json:"field"`
}

type DataRecord struct {
	Row    map[string]any           `json:"row"`
	Access []datamodel.AccessConfig `json:"access"`
}
