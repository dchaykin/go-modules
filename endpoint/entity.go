package endpoint

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/database"
	"github.com/dchaykin/go-modules/datamodel"
	"github.com/dchaykin/go-modules/helper"
	"github.com/dchaykin/go-modules/httpcomm"
	"github.com/dchaykin/go-modules/log"
	"github.com/dchaykin/go-modules/overview"
	"github.com/gorilla/mux"
)

func GetMenuItemsFromRequest(w http.ResponseWriter, r *http.Request, appName, subPath string) {
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

	mc := datamodel.MenuConfig{}
	err = mc.ReadFromFile(os.Getenv("ASSETS_PATH")+"config/"+subPath+tenant, version)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
		return
	}

	result := mc.CreateMenuByRole(userIdentity.RoleByApp(appName))

	httpcomm.ServiceResponse{
		Data: result,
	}.WriteData(w, httpcomm.PayloadFormatJSON)
}

func GetTenantConfig(w http.ResponseWriter, r *http.Request, configPath, appName string) *datamodel.TenantConfig {
	tenant, version, err := GetTenantVersionFromRequest(r)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusBadRequest)
		return nil
	}

	userIdentity, err := auth.GetUserIdentityFromRequest(*r)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusUnauthorized)
		return nil
	}

	path := fmt.Sprintf("%s/%s", configPath, tenant)
	tenantConfig, err := datamodel.LoadDataModelByRole(path, userIdentity.RoleByApp(appName), version)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
		return nil
	}

	log.Debug("Loaded tenant config for %s, version %d, app %s, subject %s", tenant, tenantConfig.Version, appName, tenantConfig.Subject)

	domainEntity := tenantConfig.DataModel[tenantConfig.Subject]
	uuid, err := datamodel.GenerateUUID()
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
		return nil
	}
	domainEntity.SetValue("uuid", uuid)

	return tenantConfig
}

func GetDomainEntityByUUID(w http.ResponseWriter, r *http.Request, domainEntity datamodel.DomainEntity) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	if uuid == "" {
		httpcomm.SetResponseError(&w, "no uuid found in the request", nil, http.StatusBadRequest)
		return
	}

	err := database.GetDomainEntityByUUID(uuid, domainEntity)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
		return
	}

	if domainEntity == nil {
		httpcomm.SetResponseError(&w, "", fmt.Errorf("no record with UUID %s found", uuid), http.StatusNotFound)
		return
	}

	httpcomm.ServiceResponse{
		Data: domainEntity.Entity(),
	}.WriteData(w, httpcomm.PayloadFormatJSON)
}

func CreateEntity(w http.ResponseWriter, r *http.Request, domainEntity datamodel.DomainEntity, subject string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		httpcomm.SetResponseError(&w, "Unable to fetch payload from body", err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Debug("create %s, body: %s", subject, string(body))

	userIdentity, err := auth.GetUserIdentityFromRequest(*r)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusUnauthorized)
		return
	}

	err = json.Unmarshal(body, domainEntity)
	if err != nil {
		httpcomm.SetResponseError(&w, "Unable to unmarshal payload", err, http.StatusBadRequest)
		return
	}

	domainEntity.SetMetadata(userIdentity, subject)

	err = saveEntity(domainEntity)
	if err != nil {
		httpcomm.SetResponseError(&w, fmt.Sprintf("Unable to save %s into the database. UUID: %s", subject, domainEntity.UUID()), err, http.StatusInternalServerError)
		return
	}

	err = overview.UpdateOverviewRow(userIdentity, domainEntity)
	if err != nil {
		httpcomm.SetResponseError(&w, fmt.Sprintf("could not create or update an overview row. UUID: %s", domainEntity.UUID()), err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func saveEntity(domainEntity datamodel.DomainEntity) error {
	domainEntity.CleanNil()
	err := helper.EnsureUUID(domainEntity)
	if err != nil {
		return fmt.Errorf("unable to generate a uuid: %v", err)
	}

	err = domainEntity.BeforeSave()
	if err != nil {
		return err
	}

	session, err := database.OpenSession()
	if err != nil {
		return log.WrapError(err)
	}
	defer session.Close()

	return session.ReplaceEntityByUUID(domainEntity, true)
}
