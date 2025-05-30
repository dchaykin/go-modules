package overview

import "github.com/dchaykin/go-modules/datamodel"

type DataRecord struct {
	Row    map[string]any           `json:"row"`
	Access []datamodel.AccessConfig `json:"access"`
}

func (r DataRecord) UUID() string {
	if r.Row == nil {
		return ""
	}
	uuid, ok := r.Row["uuid"]
	if !ok {
		return ""
	}
	return uuid.(string)
}
