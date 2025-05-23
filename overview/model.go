package overview

import "github.com/dchaykin/go-modules/datamodel"

type OverviewAction struct {
	Command string  `json:"command"`
	Icon    *string `json:"icon,omitempty"`
	Link    string  `json:"link"`
	Field   string  `json:"field"`
}

type SubjectConfig struct {
	Name       string                 `json:"name"`
	ActionList []OverviewAction       `json:"actions"`
	Fields     datamodel.CustomFields `json:"fields"`
}

type DataRecord struct {
	Row    map[string]any           `json:"row"`
	Access []datamodel.AccessConfig `json:"access"`
}

type OverviewCommand struct {
	Subject    string `json:"-"`
	ServiceUrl string `json:"-"`

	CreateTable    *bool            `json:"createTable,omitempty"`
	ActionList     []OverviewAction `json:"actions,omitempty"`
	FillData       *bool            `json:"fillData,omitempty"`
	FillComboboxes *bool            `json:"fillComboboxes,omitempty"`
	FillCombobox   *string          `json:"fillCombobox,omitempty"`
}
