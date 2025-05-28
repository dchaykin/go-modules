package datamodel

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	FieldTypeInt      = "int"
	FieldTypeUint     = "uint"
	FieldTypeString   = "string"
	FieldTypeFloat    = "float"
	FieldTypeDate     = "date"
	FieldTypeBool     = "bool"
	FieldTypeDateTime = "datetime"
	FieldTypeCombobox = "cmb"
	FieldTypeImage    = "image"
	FieldTypeList     = "list"
	FieldTypeFile     = "file"
)

type CustomFields map[string]CustomField

func (cfs *CustomFields) SetValue(fieldName string, value any) {
	if (*cfs)[fieldName] == nil {
		(*cfs)[fieldName] = CustomField{}
	}
	(*cfs)[fieldName]["value"] = value
}

type FieldConfig struct {
	Mandatory *bool   `json:"mandatory,omitempty"`
	Readonly  *bool   `json:"readonly,omitempty"`
	Masked    *bool   `json:"masked,omitempty"`
	Command   *string `json:"command,omitempty"`
}

func (fc FieldConfig) isMandatory() bool {
	if fc.Mandatory == nil {
		return false
	}
	return *fc.Mandatory
}

func (fc FieldConfig) isReadonly() bool {
	if fc.Readonly == nil {
		return false
	}
	return *fc.Readonly
}

func (fc FieldConfig) getCommand() string {
	if fc.Command == nil {
		return ""
	}
	return *fc.Command
}

type UserCommand []string

type RecordConfig map[string]FieldConfig

func (rc *RoleConfig) UnmarshalJSON(data []byte) error {
	temp := RoleConfig{}
	if err := json.Unmarshal(data, &temp.RecordConfig); err != nil {
		return err
	}
	*rc = temp
	return nil
}

type RoleConfig struct {
	RecordConfig map[string]RecordConfig `json:"roles"`
}

type ComboboxInfo struct {
	Version  int    `json:"version"`
	FileName string `json:"source"`
}

type TenantConfig struct {
	Version      int                      `json:"version"`
	Subject      string                   `json:"subject"`
	DataModel    map[string]CustomFields  `json:"datamodel"`
	Roles        *map[string]RoleConfig   `json:"roles,omitempty"`
	Layout       any                      `json:"layout"`
	ComboboxInfo *ComboboxInfo            `json:"cmbsInfo,omitempty"`
	Cmbs         *TenantComboboxDatamodel `json:"cmbs,omitempty"`
	Prefix       map[string]string        `json:"prefix"`
}

func (tc TenantConfig) GetPrefix(key string) string {
	if value, ok := tc.Prefix[key]; ok {
		return value
	}
	return ""
}

type CustomField map[string]any

func (cf CustomField) Type() string {
	if cf["type"] == nil {
		return FieldTypeString
	}
	return fmt.Sprintf("%s", cf["type"])
}

func (cf *CustomField) setMandatory(mandatory bool) {
	(*cf)["mandatory"] = nil
	if mandatory {
		(*cf)["mandatory"] = &mandatory
	}
}

func (cf *CustomField) setReadonly(readOnly bool) {
	(*cf)["readonly"] = nil
	if readOnly {
		(*cf)["readonly"] = &readOnly
	}
}

func (cf *CustomField) setCommand(command string) {
	(*cf)["command"] = nil
	if command != "" {
		(*cf)["command"] = command
	}
}

func (cf CustomField) IsMandatory() bool {
	result, ok := cf["mandatory"]
	if !ok || result == nil {
		return false
	}
	return *result.(*bool)
}

func (cf CustomField) IsReadonly() bool {
	result, ok := cf["readonly"]
	if !ok || result == nil {
		return false
	}
	return *result.(*bool)
}

func (cf CustomField) IsMasked() bool {
	result, ok := cf["masked"]
	if !ok || result == nil {
		return false
	}
	return result.(bool)
}

func loadDataModelFromFile(path string, version int) (*TenantConfig, error) {
	fileName := fmt.Sprintf("%s-%03d/datamodel.json", path, version)
	jsonData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	tc := TenantConfig{}
	if err = json.Unmarshal(jsonData, &tc); err != nil {
		return nil, err
	}

	if tc.Version != version {
		return nil, fmt.Errorf("tenant version does not match. Expected %d, got %d", version, tc.Version)
	}

	if tc.Subject == "" {
		return nil, fmt.Errorf("subject is empty")
	}

	return &tc, nil
}

func ReadPrefix(path, key string, version int) string {
	tc, err := loadDataModelFromFile(path, version)
	if err != nil {
		return ""
	}
	return tc.GetPrefix(key)
}

func LoadDataModel(path string, version int) (*TenantConfig, error) {
	tc, err := loadDataModelFromFile(path, version)
	if err != nil {
		return nil, err
	}

	if tc.ComboboxInfo != nil && tc.ComboboxInfo.FileName != "" {
		tenantCmbs, err := loadTenantComboboxList(path, tc.ComboboxInfo.FileName, tc.ComboboxInfo.Version)
		if err != nil {
			return nil, err
		}
		tc.Cmbs = tenantCmbs
	}

	tc.ComboboxInfo = nil

	return tc, nil
}

func LoadDataModelByRole(path, role string, version int) (*TenantConfig, error) {
	tc, err := LoadDataModel(path, version)
	if err != nil {
		return nil, err
	}

	if tc.Roles == nil {
		return tc, nil
	}

	role = strings.ToLower(role)
	roleConfig, ok := (*tc.Roles)[role]
	if !ok {
		return nil, fmt.Errorf("no config for role %s found", role)
	}

	for recordName, recordConfig := range roleConfig.RecordConfig {
		record, ok := tc.DataModel[recordName]
		if !ok {
			return nil, fmt.Errorf("record config %s for role %s exists, but no record was found in datamodel", recordName, role)
		}
		for fieldName, fieldConfig := range recordConfig {
			field, ok := record[fieldName]
			if !ok {
				return nil, fmt.Errorf("field config %s.%s for role %s exists, but no record was found in datamodel", recordName, fieldName, role)
			}
			field.setMandatory(fieldConfig.isMandatory())
			field.setReadonly(fieldConfig.isReadonly())
			field.setCommand(fieldConfig.getCommand())
		}
	}

	tc.Roles = nil

	return tc, nil
}
