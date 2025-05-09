package datamodel

import (
	"encoding/json"
	"fmt"
	"os"
)

type Combobox struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type TenantCombobox struct {
	Translate *bool      `json:"translate,omitempty"`
	Content   []Combobox `json:"content"`
	Source    *string    `json:"source,omitempty"`
	Name      string     `json:"name"`
	Type      *string    `json:"type,omitempty"`
}

type TenantComboboxList map[string]TenantCombobox
type TenantComboboxDatamodel map[string]TenantComboboxList

func loadTenantComboboxList(path2config, cmbFileName string, version int) (*TenantComboboxDatamodel, error) {
	fileName := fmt.Sprintf("%s/%s", path2config, cmbFileName)
	jsonData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	tcd := TenantComboboxDatamodel{}
	if err = json.Unmarshal(jsonData, &tcd); err != nil {
		return nil, err
	}
	return &tcd, nil
}
