package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/datamodel"
	"github.com/dchaykin/go-modules/httpcomm"
	"github.com/dchaykin/go-modules/log"
	"github.com/dchaykin/go-modules/overview"
)

func CreateTemporaryOverview(w http.ResponseWriter, r *http.Request, subPath string) {
	tenant, version, err := GetTenantVersionFromRequest(r)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusBadRequest)
		return
	}

	userIdentity, err := auth.GetUserIdentityFromRequest(*r)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusUnauthorized)
		return
	}

	if subPath != "" {
		subPath += "/"
	}

	configPath := os.Getenv("ASSETS_PATH") + "config/" + subPath + tenant

	log.Info("Creating overview for tenant %s, role 'default', version %d", tenant, version)

	tc, err := datamodel.LoadDataModelByRole(configPath, "default", version)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
		return
	}

	err = overview.ConfigureOverview(userIdentity, *tc, tenant)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
		return
	}

	httpcomm.ServiceResponse{}.WriteData(w, httpcomm.PayloadFormatJSON)
}

func BulkInsertIntoOverview(w http.ResponseWriter, r *http.Request, subject string, entityList []datamodel.DomainEntity, isTemporary bool) {
	userIdentity, err := auth.GetUserIdentityFromRequest(*r)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusUnauthorized)
		return
	}

	recordList := []overview.DataRecord{}
	for _, entity := range entityList {
		record := overview.DataRecord{
			Row:    entity.OverviewRow(),
			Access: entity.GetAccessConfig(),
		}
		recordList = append(recordList, record)
	}

	log.Info("Bulk inserting %d records into overview for subject '%s'", len(recordList), subject)

	payload, err := json.Marshal(recordList)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
		return
	}

	params := map[string]string{}
	if isTemporary {
		params["temporary"] = "true"
	}

	endpoint := fmt.Sprintf("https://%s/app-overview/api/bulk-insert/%s", os.Getenv("MYHOST"), subject)
	resp := httpcomm.Post(endpoint, userIdentity, params, string(payload))
	if resp.StatusCode != http.StatusOK {
		httpcomm.SetResponseError(&w, "", resp.GetError(), http.StatusInternalServerError)
		return
	}

	log.Info("Bulk insert into overview for subject '%s' completed successfully", subject)

	httpcomm.ServiceResponse{}.WriteData(w, httpcomm.PayloadFormatJSON)
}

func CommitOverview(w http.ResponseWriter, r *http.Request, tenant, subject string) {

	log.Info("Committing overview for tenant '%s', subject '%s'", tenant, subject)

	userIdentity, err := auth.GetUserIdentityFromRequest(*r)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusUnauthorized)
		return
	}

	endpoint := fmt.Sprintf("https://%s/app-overview/api/commit/overview/%s/%s", os.Getenv("MYHOST"), tenant, subject)
	resp := httpcomm.Post(endpoint, userIdentity, nil, "")
	if resp.StatusCode != http.StatusOK {
		httpcomm.SetResponseError(&w, "", resp.GetError(), http.StatusInternalServerError)
		return
	}

	log.Info("Overview for tenant '%s', subject '%s' committed successfully", tenant, subject)

	httpcomm.ServiceResponse{}.WriteData(w, httpcomm.PayloadFormatJSON)
}
