package overview

import "github.com/dchaykin/go-modules/datamodel"

type DataRecord struct {
	Row    map[string]any           `json:"row"`
	Access []datamodel.AccessConfig `json:"access"`
}
