package overview

import "github.com/dchaykin/go-modules/datamodel"

type SubjectConfig struct {
	ServiceURL string                 `json:"serviceUrl"`
	Name       string                 `json:"name"`
	Fields     datamodel.CustomFields `json:"fields"`
}

type DataRecord struct {
	Row    map[string]interface{}   `json:"row"`
	Access []datamodel.AccessConfig `json:"access"`
}

type OverviewCommand struct {
	Subject    string `json:"-"`
	ServiceUrl string `json:"-"`

	CreateTable    *bool   `json:"createTable,omitempty"`
	FillData       *bool   `json:"fillData,omitempty"`
	FillComboboxes *bool   `json:"fillComboboxes,omitempty"`
	FillCombobox   *string `json:"fillCombobox,omitempty"`
}
