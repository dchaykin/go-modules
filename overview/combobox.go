package overview

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/datamodel"
	"github.com/dchaykin/go-modules/httpcomm"
	"github.com/dchaykin/go-modules/log"
)

func CreateComboboxes(userIdentity auth.UserIdentity, path, overviewName string) error {
	tc, err := datamodel.LoadDataModel(path)
	if err != nil {
		return log.WrapError(err)
	}

	if tc.Cmbs == nil {
		log.Info("No comboboxes for %s found", overviewName)
		return nil
	}

	for fieldName, customField := range tc.DataModel[overviewName] {
		if customField.Type() == datamodel.FieldTypeCombobox {
			rootNode, ok := (*tc.Cmbs)[overviewName]
			if !ok {
				continue
			}
			for _, cmb := range rootNode {
				if err = CreateCombobox(userIdentity, overviewName, fieldName, cmb); err != nil {
					return log.WrapError(err)
				}
			}
		}
	}

	return nil
}

func CreateCombobox(userIdentity auth.UserIdentity, overviewName, fieldName string, cmb datamodel.TenantCombobox) error {
	data, err := json.Marshal(cmb)
	if err != nil {
		return log.WrapError(err)
	}
	endpoint := fmt.Sprintf("https://%s/app-overview/api/create/combobox/%s/%s", os.Getenv("MYHOST"), overviewName, fieldName)
	resp := httpcomm.Post(endpoint, userIdentity, nil, string(data))
	if resp.StatusCode != http.StatusOK {
		return log.WrapError(resp.GetError())
	}
	log.Info("Combobox %s (re)created", fieldName)
	return nil
}
