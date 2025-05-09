package database

import (
	"fmt"
	"net/http"

	"github.com/dchaykin/go-modules/datamodel"
	"github.com/dchaykin/go-modules/httpcomm"
	"github.com/gorilla/mux"
)

func GetDomainEntityByUUID(r *http.Request, domainEntity datamodel.DomainEntity) (*httpcomm.ServiceResponse, int) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	if uuid == "" {
		return datamodel.GetErrorResponse(fmt.Errorf("no uuid found in the request")), http.StatusBadRequest
	}

	err := loadDomainEntityByUUID(uuid, domainEntity)
	if err != nil {
		return datamodel.GetErrorResponse(err), http.StatusInternalServerError
	}

	if domainEntity == nil {
		return datamodel.GetErrorResponse(fmt.Errorf("no record with UUID %s found", uuid)), http.StatusNotFound
	}

	return &httpcomm.ServiceResponse{
		Data: domainEntity.Entity(),
	}, http.StatusOK
}

func loadDomainEntityByUUID(uuid string, domainEntity datamodel.DomainEntity) error {
	session, err := OpenSession()
	if err != nil {
		return err
	}
	defer session.Close()

	bFound, err := session.GetEntityByUUID(uuid, domainEntity)
	if err != nil {
		return err
	}
	if !bFound {
		return fmt.Errorf("no record with UUID %s found", uuid)
	}
	return nil
}
