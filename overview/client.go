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

func CreateOverviewConfig(userIdentity auth.UserIdentity, path, serviceURL, name string) (string, error) {
	tc, err := datamodel.LoadDataModel(path)
	if err != nil {
		return "", err
	}

	subjectConfig := SubjectConfig{
		ServiceURL: serviceURL,
		Name:       name,
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
