package overview

import (
	"net/http"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/datamodel"
	"github.com/dchaykin/go-modules/httpcomm"
)

func CreateOverview(w http.ResponseWriter, r *http.Request, appName, configPath string) {
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
