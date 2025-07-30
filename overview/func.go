package overview

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/database"
	"github.com/dchaykin/go-modules/datamodel"
	"github.com/dchaykin/go-modules/httpcomm"
	"github.com/dchaykin/go-modules/log"
)

func CreateTemporaryOverview(userIdentity auth.UserIdentity, pathToDatamodel string) error {
	tenant := userIdentity.Tenant()

	log.Info("Creating overview for datamodel %s", pathToDatamodel)

	tenantConfig, err := datamodel.LoadDataModelByRole(pathToDatamodel, "default")
	if err != nil {
		return err
	}

	payload, err := json.Marshal(tenantConfig)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("https://%s/app-overview/api/create/overview/%s?temporary=true", os.Getenv("MYHOST"), tenant)
	resp := httpcomm.Post(endpoint, userIdentity, nil, payload)
	if resp.StatusCode != http.StatusOK {
		return resp.GetError()
	}

	log.Info("Overview %s config for tenant %s created, version %d. Response: %s", tenantConfig.Subject, tenant, tenantConfig.Version, string(resp.Answer))
	return nil
}

func BulkInsertIntoOverview(userIdentity auth.UserIdentity, subject string, entityList []database.DomainEntity, isTemporary bool) error {
	recordList := []DataRecord{}
	for _, entity := range entityList {

		entity.NormalizePrimitives()
		entity.ApplyMapper()

		record := DataRecord{
			Row:    entity.OverviewRow(),
			Access: entity.GetAccessConfig(),
		}
		recordList = append(recordList, record)
	}

	log.Info("Bulk inserting %d records into overview for subject '%s'", len(recordList), subject)

	payload, err := json.Marshal(recordList)
	if err != nil {
		return err
	}

	params := ""
	if isTemporary {
		params = "temporary=true"
	}

	endpoint := fmt.Sprintf("https://%s/app-overview/api/bulk-insert/%s?%s", os.Getenv("MYHOST"), subject, params)
	resp := httpcomm.Post(endpoint, userIdentity, nil, payload)
	if resp.StatusCode != http.StatusOK {
		return resp.GetError()
	}

	log.Info("Bulk insert into overview for subject '%s' completed successfully", subject)
	return nil
}

func CommitOverview(userIdentity auth.UserIdentity, subject string) error {
	tenant := userIdentity.Tenant()

	log.Info("Committing overview for tenant '%s', subject '%s'", tenant, subject)

	endpoint := fmt.Sprintf("https://%s/app-overview/api/commit/overview/%s/%s", os.Getenv("MYHOST"), tenant, subject)
	resp := httpcomm.Post(endpoint, userIdentity, nil, nil)
	if resp.StatusCode != http.StatusOK {
		return resp.GetError()
	}

	log.Info("Overview for tenant '%s', subject '%s' committed successfully", tenant, subject)
	return nil
}

func UpdateOverviewRow(domainEntity database.DomainEntity) error {
	domainEntity.ApplyMapper()

	data := DataRecord{
		Row:    domainEntity.OverviewRow(),
		Access: domainEntity.GetAccessConfig(),
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return log.WrapError(err)
	}

	endpoint := fmt.Sprintf("https://%s/app-overview/api/save/%s", os.Getenv("MYHOST"), domainEntity.CollectionName())
	resp := httpcomm.Post(endpoint, domainEntity.UserIdentity(), nil, payload)
	if resp.StatusCode != http.StatusOK {
		return resp.GetError()
	}

	return nil
}
