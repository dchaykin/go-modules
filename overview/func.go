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

func ConfigureOverview(userIdentity auth.UserIdentity, tenantConfig datamodel.TenantConfig, tenant string) error {
	payload, err := json.Marshal(tenantConfig)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("https://%s/app-overview/api/create/overview/%s?temporary=true", os.Getenv("MYHOST"), tenant)
	resp := httpcomm.Post(endpoint, userIdentity, nil, string(payload))
	if resp.StatusCode != http.StatusOK {
		return resp.GetError()
	}

	log.Info("Overview %s config for tenant %s created, version %d. Response: %s", tenantConfig.Subject, tenant, tenantConfig.Version, string(resp.Answer))
	return nil
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

	endpoint := fmt.Sprintf("https://%s/app-overview/api/save/%s", os.Getenv("MYHOST"), domainEntity.CollectionName())
	resp := httpcomm.Post(endpoint, userIdentity, nil, string(payload))
	if resp.StatusCode != http.StatusOK {
		return resp.GetError()
	}

	return nil
}
