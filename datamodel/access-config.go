package datamodel

type AccessConfig struct {
	Partner   string  `json:"partner"`
	Algorithm *string `json:"algo,omitempty"`
}
