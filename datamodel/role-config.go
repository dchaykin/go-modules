package datamodel

import (
	"encoding/json"
	"fmt"
	"os"
)

/***********************************************/
/*                 Role Files                  */
/***********************************************/

type roleFiles struct {
	ComboboxFile string `json:"combobox"`
	OverviewFile string `json:"overview"`
	FieldFile    string `json:"field"`
}

func (rf roleFiles) getComboboxes(path2config string) (*TenantComboboxDatamodel, error) {
	if rf.ComboboxFile == "" {
		return nil, nil
	}

	fileName := fmt.Sprintf("%s/%s", path2config, rf.ComboboxFile)
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

func (rf roleFiles) getFields(path2config string) (map[string]recordConfig, error) {
	if rf.FieldFile == "" {
		return nil, nil
	}

	fileName := fmt.Sprintf("%s/%s", path2config, rf.FieldFile)
	jsonData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	fcfg := map[string]recordConfig{}
	if err = json.Unmarshal(jsonData, &fcfg); err != nil {
		return nil, err
	}
	return fcfg, nil
}

func (rf roleFiles) getOverviews(path2config string) (TenantOverviewDatamodel, error) {
	if rf.OverviewFile == "" {
		return nil, nil
	}

	fileName := fmt.Sprintf("%s/%s", path2config, rf.OverviewFile)
	jsonData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	ocl := map[string]overviewSubject{}
	if err = json.Unmarshal(jsonData, &ocl); err != nil {
		return nil, err
	}
	return ocl, nil
}

/***********************************************/
/*                   Fields                    */
/***********************************************/

type fieldConfig struct {
	Mandatory *bool   `json:"mandatory,omitempty"`
	Readonly  *bool   `json:"readonly,omitempty"`
	Masked    *bool   `json:"masked,omitempty"`
	Command   *string `json:"command,omitempty"`
}

func (fc fieldConfig) isMandatory() bool {
	if fc.Mandatory == nil {
		return false
	}
	return *fc.Mandatory
}

func (fc fieldConfig) isReadonly() bool {
	if fc.Readonly == nil {
		return false
	}
	return *fc.Readonly
}

func (fc fieldConfig) getCommand() string {
	if fc.Command == nil {
		return ""
	}
	return *fc.Command
}

type recordConfig map[string]fieldConfig

/***********************************************/
/*                 Comboboxes                  */
/***********************************************/

type ComboboxType string

const (
	ComboboxTypeUnknown = ""
	ComboboxTypeStatic  = "static"
	ComboboxTypeApi     = "api"
	ComboboxTypeSelf    = "self"
)

type Combobox struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type TenantCombobox struct {
	Translate *bool         `json:"translate,omitempty"`
	Content   []Combobox    `json:"content"`
	Source    *string       `json:"source,omitempty"`
	Name      string        `json:"name"`
	Type      *ComboboxType `json:"type,omitempty"`
}

func (tc TenantCombobox) GetType() ComboboxType {
	if tc.Type == nil {
		return ComboboxTypeStatic
	}
	return *tc.Type
}

type TenantComboboxList map[string]TenantCombobox
type TenantComboboxDatamodel map[string]TenantComboboxList

/***********************************************/
/*              Overview Configs               */
/***********************************************/

type overviewConfig struct {
	Name string `json:"name"` // Unique Key
}

type overviewSubject struct {
	CommandList  []overviewCommand `json:"command"`
	OverviewList []overviewConfig  `json:"overview"`
}

type TenantOverviewDatamodel map[string]overviewSubject

func (ov *overviewSubject) mergeOverview(source overviewSubject) {
	ov.mergeCommandList(source.CommandList)
	ov.mergeOverviewList(source.OverviewList)
}

func (ov *overviewSubject) mergeCommandList(srcCommandList []overviewCommand) {
	for _, sourceCmd := range srcCommandList {
		cmd := ov.getCommandByAction(sourceCmd.Action)
		if cmd == nil {
			ov.CommandList = append(ov.CommandList, sourceCmd)
			continue
		}
		*cmd = sourceCmd
	}
}

func (ov overviewSubject) getCommandByAction(action string) *overviewCommand {
	for i, cmd := range ov.CommandList {
		if cmd.Action == action {
			return &ov.CommandList[i]
		}
	}
	return nil
}

func (ov *overviewSubject) mergeOverviewList(srcOverviewList []overviewConfig) {
	for _, sourceOverview := range srcOverviewList {
		overviewConfig := ov.getOverviewByName(sourceOverview.Name)
		if overviewConfig == nil {
			ov.OverviewList = append(ov.OverviewList, sourceOverview)
			continue
		}
		*overviewConfig = sourceOverview
	}
}

func (ov overviewSubject) getOverviewByName(name string) *overviewConfig {
	for i, overview := range ov.OverviewList {
		if overview.Name == name {
			return &ov.OverviewList[i]
		}
	}
	return nil
}
