package httpcomm

import (
	"encoding/json"
	"fmt"
	"time"
)

type MetaData struct {
	OriginalFileName string    `json:"filename"`
	FileSize         int64     `json:"size"`
	Author           string    `json:"author"`
	Partner          string    `json:"partner"`
	MimeType         string    `json:"mime"`
	Hash             *string   `json:"hash,omitempty"`
	Error            *string   `json:"error,omitempty"`
	Timestamp        time.Time `json:"createdAt"`
}

func (md MetaData) Stringer() string {
	data, err := json.Marshal(md)
	if err != nil {
		return fmt.Sprintf(`{"error":"%v"}`, err)
	}
	return string(data)
}

func (md *MetaData) Set(data []byte) error {
	return json.Unmarshal(data, md)
}

func (md *MetaData) SetHash(hashValue string) {
	md.Hash = &hashValue
}

func (md MetaData) GetHash() string {
	if md.Hash == nil {
		return ""
	}
	return *md.Hash
}
