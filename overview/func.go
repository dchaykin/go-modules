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

func configureOverview(userIdentity auth.UserIdentity, tenantConfig datamodel.TenantConfig, tenant, userRole string) error {
	payload, err := json.Marshal(tenantConfig)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("https://%s/app-overview/api/create/overview/%s/%s", os.Getenv("MYHOST"), tenant, userRole)
	resp := httpcomm.Post(endpoint, userIdentity, nil, string(payload))
	if resp.StatusCode != http.StatusOK {
		return resp.GetError()
	}

	log.Info("Overview %s config for tenant %s and role %s created, version %d. Response: %s", tenantConfig.Subject, tenant, userRole, tenantConfig.Version, string(resp.Answer))
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

func CreateOverview(w http.ResponseWriter, r *http.Request, appName, subPath string) {
	tenant, version, err := httpcomm.GetTenantVersionFromRequest(r)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusBadRequest)
		return
	}

	userIdentity, err := auth.GetUserIdentityFromRequest(*r)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusUnauthorized)
		return
	}

	roleName := userIdentity.RoleByApp(appName)

	if subPath != "" {
		subPath += "/"
	}

	configPath := "config/" + subPath + tenant

	tc, err := datamodel.LoadDataModelByRole(configPath, roleName, version)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
		return
	}

	err = configureOverview(userIdentity, *tc, tenant, roleName)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
		return
	}

	httpcomm.ServiceResponse{}.WriteData(w, httpcomm.PayloadFormatJSON)
}
