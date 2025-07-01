package endpoint

import (
	"net/http"
	"os"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/datamodel"
	"github.com/dchaykin/go-modules/httpcomm"
	"github.com/dchaykin/go-modules/overview"
)

func CreateOverview(w http.ResponseWriter, r *http.Request, appName, subPath string) {
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

	roleName := userIdentity.RoleByApp(appName)

	if subPath != "" {
		subPath += "/"
	}

	configPath := os.Getenv("ASSETS_PATH") + "config/" + subPath + tenant

	tc, err := datamodel.LoadDataModelByRole(configPath, roleName, version)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
		return
	}

	err = overview.ConfigureOverview(userIdentity, *tc, tenant, roleName)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
		return
	}

	httpcomm.ServiceResponse{}.WriteData(w, httpcomm.PayloadFormatJSON)
}
