package datamodel

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/httpcomm"
)

type Metadata struct {
	Timestamp time.Time `bson:"timestamp"`
	User      string    `bson:"user"`
	Partner   string    `bson:"partner"`
	Role      string    `bson:"role"`
}

type Record struct {
	Metadata Metadata       `bson:"metadata"`
	Fields   map[string]any `json:"entity" bson:"entity"`
	Mapper   map[string]any `json:"mapper" bson:"mapper"`
}

func (r *Record) SetMetadata(userIdentity auth.UserIdentity, subject string) {
	r.Metadata.Timestamp = time.Now()
	r.Metadata.Partner = userIdentity.Partner()
	r.Metadata.Role = userIdentity.RoleBySubject(subject)
	r.Metadata.User = userIdentity.Username()
}

type OnJsonArrayFound func(array []any)

func (r *Record) FindJsonArray(jsonPath []string, f OnJsonArrayFound) {
	r.findJsonArray(r.Fields, jsonPath, f)
}

func (r *Record) findJsonArray(node map[string]any, jsonPath []string, f OnJsonArrayFound) {
	if len(jsonPath) == 0 {
		return
	}
	lookedNodeName := jsonPath[0]
	for k := range node {
		if k == lookedNodeName {
			switch vTyped := node[k].(type) {
			case []any:
				f(vTyped)
				return
			case map[string]any:
				r.findJsonArray(vTyped, jsonPath[1:], f)
			}
		}
	}
}

type OnJsonFieldFound func(field map[string]any, name string)

func (r *Record) FindJsonField(jsonPath []string, f OnJsonFieldFound) {
	r.findJsonField(r.Fields, jsonPath, f)
}

func (r *Record) findJsonField(node map[string]any, jsonPath []string, f OnJsonFieldFound) {
	if len(jsonPath) == 0 {
		return
	}
	lookedNodeName := jsonPath[0]
	for k, v := range node {
		if k == lookedNodeName {
			switch vTyped := v.(type) {
			case nil:
				if len(jsonPath) == 1 {
					f(node, k)
				}
			case []any:
				for i := range vTyped {
					if vTyped[i] == nil {
						continue
					}
					r.findJsonField(vTyped[i].(map[string]any), jsonPath[1:], f)
				}
				return
			case map[string]any:
				r.findJsonField(node[k].(map[string]any), jsonPath[1:], f)
			case any:
				if len(jsonPath) == 1 {
					f(node, k)
				}
			}
		}
	}
}

func (r Record) UUID() string {
	if uuid, ok := r.Fields["uuid"]; ok {
		return fmt.Sprintf("%s", uuid)
	}
	return ""
}

func (r *Record) SetUUID(UUID string) {
	r.Fields["uuid"] = UUID
}

func (r Record) Entity() map[string]any {
	return r.Fields
}

func GetErrorResponse(err error) *httpcomm.ServiceResponse {
	result := httpcomm.ServiceResponse{Error: new(string)}
	*result.Error = fmt.Sprintf("%v", err)
	return &result
}

func GetDomainConfig(r *http.Request, configPath, subject string) (*httpcomm.ServiceResponse, int) {
	tenant, version, err := httpcomm.GetTenantVersionFromRequest(r)
	if err != nil {
		return GetErrorResponse(err), http.StatusBadRequest
	}

	userIdentity, err := auth.GetUserIdentityFromRequest(*r)
	if err != nil {
		return GetErrorResponse(err), http.StatusUnauthorized
	}

	path := fmt.Sprintf("%s/%s", configPath, tenant)
	tenantConfig, err := LoadDataModelByRole(path, userIdentity.RoleBySubject(subject), version)
	if err != nil {
		return GetErrorResponse(err), http.StatusInternalServerError
	}

	domainEntity := tenantConfig.DataModel[subject]
	uuid, err := GenerateUUID()
	if err != nil {
		return GetErrorResponse(err), http.StatusInternalServerError
	}
	domainEntity.SetValue("uuid", uuid)

	return &httpcomm.ServiceResponse{Data: tenantConfig}, http.StatusOK
}
