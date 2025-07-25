package datamodel

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/dchaykin/go-modules/log"
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

type TenantConfig struct {
	path      string
	Version   int                      `json:"version"`
	Subject   string                   `json:"subject"`
	DataModel map[string]CustomFields  `json:"datamodel"`
	Roles     *map[string]roleFiles    `json:"roles,omitempty"`
	Layout    any                      `json:"layout"`
	Cmbs      *TenantComboboxDatamodel `json:"cmbs,omitempty"`
	Overviews *OverviewModel           `json:"overview,omitempty"`
	Prefix    map[string]string        `json:"prefix"`
}

func (tc TenantConfig) GetPrefix(key string) string {
	if value, ok := tc.Prefix[key]; ok {
		return value
	}
	return ""
}

func (tc *TenantConfig) setReadonly(readonly bool) {
	for recordName := range tc.DataModel {
		record, ok := tc.DataModel[recordName]
		if !ok {
			continue
		}
		for fieldName := range record {
			field, ok := record[fieldName]
			if !ok {
				continue
			}
			field.setReadonly(readonly)
		}
	}
}

func (tc *TenantConfig) buildRole(roleConfig roleFiles, roleName string, isReadonlyIfMissing bool) error {
	// Comboboxes
	cmbs, err := roleConfig.getComboboxes(tc.path)
	if err != nil {
		return err
	}

	if cmbs != nil {
		if tc.Cmbs == nil {
			tc.Cmbs = cmbs
		} else {
			for recordName, recordConfig := range *cmbs {
				(*tc.Cmbs)[recordName] = recordConfig
			}
		}
	} else if tc.Cmbs == nil {
		tc.Cmbs = &TenantComboboxDatamodel{}
	}

	// Fields
	fields, err := roleConfig.getFields(tc.path)
	if err != nil {
		return err
	}

	if fields == nil {
		tc.setReadonly(isReadonlyIfMissing)
	} else {
		for recordName, recordConfig := range fields {
			record, ok := tc.DataModel[recordName]
			if !ok {
				return fmt.Errorf("record config %s for role %s exists, but no record was found in datamodel", recordName, roleName)
			}
			for fieldName, fieldConfig := range recordConfig {
				field, ok := record[fieldName]
				if !ok {
					return fmt.Errorf("field config %s.%s for role %s exists, but no record was found in datamodel", recordName, fieldName, roleName)
				}
				field.setMandatory(fieldConfig.isMandatory())
				field.setReadonly(fieldConfig.isReadonly())
				field.setCommand(fieldConfig.getCommand())
			}
		}
	}

	// Overviews
	overviews, err := roleConfig.getOverviewModel(tc.path)
	if err != nil {
		return err
	}

	if overviews != nil {
		if tc.Overviews == nil {
			tc.Overviews = overviews
		} else {
			tc.Overviews = &OverviewModel{}
			tc.Overviews.mergeOverviews(*overviews)
		}
	} else if tc.Overviews == nil {
		tc.Overviews = &OverviewModel{}
	}

	return nil
}

type CustomField map[string]any

func (cf CustomField) Type() string {
	if cf["type"] == nil {
		return FieldTypeString
	}
	return fmt.Sprintf("%s", cf["type"])
}

func (cf CustomField) ValueByType(value any) any {
	if value == nil {
		return nil
	}
	switch cf.Type() {
	case FieldTypeDate:
		date := fmt.Sprintf("%s", value)
		if len(date) >= 10 {
			return date[:10]
		}
	}
	return value
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

func loadDataModelFromFile(path string) (*TenantConfig, error) {
	jsonData, err := os.ReadFile(path + "/datamodel.json")
	if err != nil {
		return nil, err
	}

	tc := TenantConfig{}
	if err = json.Unmarshal(jsonData, &tc); err != nil {
		return nil, err
	}

	if tc.Subject == "" {
		return nil, fmt.Errorf("subject is empty")
	}

	tc.path = path

	return &tc, nil
}

func ReadPrefix(path, key string) string {
	tc, err := loadDataModelFromFile(path)
	if err != nil {
		return ""
	}
	return tc.GetPrefix(key)
}

func GetRoles(path string) ([]string, error) {
	tc, err := loadDataModelFromFile(path)
	if err != nil {
		return nil, err
	}

	if tc.Roles == nil {
		log.Warn("No roles found in the data model at path %s, using default role", path)
		return []string{"default"}, nil
	}

	roles := make([]string, 0, len(*tc.Roles))
	for roleName := range *tc.Roles {
		roles = append(roles, roleName)
	}

	return roles, nil
}

func LoadDataModelByRole(fileName, roleName string) (*TenantConfig, error) {
	tc, err := loadDataModelFromFile(os.Getenv("ASSETS_PATH") + fileName)
	if err != nil {
		return nil, err
	}

	if tc.Roles == nil {
		return tc, nil
	}

	defaultConfig, ok := (*tc.Roles)["default"]
	if !ok {
		return nil, fmt.Errorf("no default config found")
	}

	err = tc.buildRole(defaultConfig, "default", false)
	if err != nil {
		return nil, err
	}

	roleName = strings.ToLower(roleName)
	roleConfig, ok := (*tc.Roles)[roleName]
	if ok {
		err = tc.buildRole(roleConfig, roleName, true)
		if err != nil {
			return nil, err
		}
	}

	tc.Roles = nil

	return tc, nil
}
