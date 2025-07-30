package overview

import "github.com/dchaykin/go-modules/database"

type DataRecord struct {
	Row    map[string]any          `json:"row"`
	Access []database.AccessConfig `json:"access"`
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
