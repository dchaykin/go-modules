package datamodel

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	FieldTypeInt      = "int"
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

func (cfs *CustomFields) SetValue(fieldName string, value interface{}) {
	if (*cfs)[fieldName] == nil {
		(*cfs)[fieldName] = CustomField{}
	}
	(*cfs)[fieldName]["value"] = value
}

type FieldConfig struct {
	Mandatory *bool   `json:"mandatory,omitempty"`
	Readonly  *bool   `json:"readonly,omitempty"`
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

type CustomField map[string]interface{}

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

func loadDataModel(path string) (*TenantConfig, error) {
	jsonData, err := readDataModel(path)
	if err != nil {
		return nil, err
	}

	tc := TenantConfig{}
	if err = json.Unmarshal(jsonData, &tc); err != nil {
		return nil, err
	}
	return &tc, nil
}

func ReadPrefix(path, key string) string {
	tc, err := loadDataModel(path)
	if err != nil {
		return ""
	}
	return tc.GetPrefix(key)
}

func LoadDataModel(path string, role string) (*TenantConfig, error) {
	tc, err := loadDataModel(path)
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

	if tc.ComboboxInfo != nil {
		tenantCmbs, err := loadTenantComboboxList(path, tc.ComboboxInfo.FileName, tc.ComboboxInfo.Version)
		if err != nil {
			return nil, err
		}
		tc.Cmbs = tenantCmbs
	}

	tc.ComboboxInfo = nil

	return tc, nil
}

func readDataModel(path string) ([]byte, error) {
	fileName := fmt.Sprintf("%s/custom-fields.json", path)
	return os.ReadFile(fileName)
}

func GetConfigPath(basePath, tenant string, version int) string {
	return fmt.Sprintf("%s/%s-%03d", basePath, tenant, version)
}
