package endpoint

import (
	"net/http"
	"strconv"
)

func GetTenantVersionFromRequest(r *http.Request) (string, int, error) {
	tenant := r.URL.Query().Get("tenant")
	versionParam := r.URL.Query().Get("version")
	version, err := strconv.Atoi(versionParam)
	return tenant, version, err
}
