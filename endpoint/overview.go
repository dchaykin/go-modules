package endpoint

import (
	"net/http"
	"os"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/datamodel"
	"github.com/dchaykin/go-modules/httpcomm"
	"github.com/dchaykin/go-modules/overview"
)

func CreateOverview(w http.ResponseWriter, r *http.Request, subPath string) {
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
	roles, err := datamodel.GetRoles(configPath, 1)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
		return
	}

	for _, role := range roles {
		tc, err := datamodel.LoadDataModelByRole(configPath, role, version)
		if err != nil {
			httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
			return
		}
		err = overview.ConfigureOverview(userIdentity, *tc, tenant, role)
		if err != nil {
			httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
			return
		}
	}

	httpcomm.ServiceResponse{}.WriteData(w, httpcomm.PayloadFormatJSON)
}
