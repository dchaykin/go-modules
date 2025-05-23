package overview

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/datamodel"
	"github.com/dchaykin/go-modules/httpcomm"
	"github.com/dchaykin/go-modules/log"
)

func createOverviewConfig(userIdentity auth.UserIdentity, path, name string, actions []OverviewAction) (string, error) {
	tc, err := datamodel.LoadDataModel(path)
	if err != nil {
		return "", err
	}

	subjectConfig := SubjectConfig{
		Name:       name,
		ActionList: actions,
		Fields:     map[string]datamodel.CustomField{},
	}
	for fieldName, customField := range tc.DataModel[name] {
		subjectConfig.Fields[fieldName] = customField
	}

	payload, err := json.Marshal(subjectConfig)
	if err != nil {
		return "", err
	}

	endpoint := fmt.Sprintf("https://%s/app-overview/api/create/table", os.Getenv("MYHOST"))
	resp := httpcomm.Post(endpoint, userIdentity, nil, string(payload))
	if resp.StatusCode != http.StatusOK {
		return "", resp.GetError()
	}

	return string(resp.Answer), nil
}

func UpdateOverviewRow(userIdentity auth.UserIdentity, domainEntity datamodel.DomainEntity) error {
	data := DataRecord{
		Row:    domainEntity.OverviewRow(),
		Access: domainEntity.GetAccessConfig(),
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return log.WrapError(err)
	}

	endpoint := fmt.Sprintf("https://%s/app-overview/api/data/%s", os.Getenv("MYHOST"), domainEntity.CollectionName())
	resp := httpcomm.Post(endpoint, userIdentity, nil, string(payload))
	if resp.StatusCode != http.StatusOK {
		return resp.GetError()
	}

	return nil
}

func PrepareOverviewCommand(w http.ResponseWriter, r *http.Request) (auth.UserIdentity, *OverviewCommand) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		httpcomm.SetResponseError(&w, "Unable to fetch payload from body", err, http.StatusBadRequest)
		return nil, nil
	}
	defer r.Body.Close()

	userIdentity, err := auth.GetUserIdentityFromRequest(*r)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusUnauthorized)
		return nil, nil
	}

	if !userIdentity.IsAdmin() {
		httpcomm.SetResponseError(&w, "permission denied", nil, http.StatusForbidden)
		return nil, nil
	}

	overviewCommand := OverviewCommand{}
	err = json.Unmarshal(body, &overviewCommand)
	if err != nil {
		httpcomm.SetResponseError(&w, "Unable to unmarshal payload", err, http.StatusBadRequest)
		return nil, nil
	}

	return userIdentity, &overviewCommand
}

func PerformOverviewCommand(userIdentity auth.UserIdentity, overviewCommand OverviewCommand, configPath string) (string, error) {
	if overviewCommand.CreateTable != nil && *(overviewCommand.CreateTable) {
		return createOverviewConfig(userIdentity, configPath, overviewCommand.Subject, overviewCommand.ActionList)
	}
	if overviewCommand.FillComboboxes != nil && *(overviewCommand.FillComboboxes) {
		err := createComboboxes(userIdentity, configPath, overviewCommand.Subject)
		if err != nil {
			return "", err
		}
	}

	return "no action", nil
}
