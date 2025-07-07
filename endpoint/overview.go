package endpoint

import (
	"net/http"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/database"
	"github.com/dchaykin/go-modules/datamodel"
	"github.com/dchaykin/go-modules/httpcomm"
	"github.com/dchaykin/go-modules/overview"
)

type OnNextBulkInsert func(session database.DatabaseSession) ([]datamodel.DomainEntity, error)

func RebuildOverview(w http.ResponseWriter, r *http.Request, subject string, f OnNextBulkInsert) {
	userIdentity, err := auth.GetUserIdentityFromRequest(*r)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusUnauthorized)
		return
	}

	_, version, err := GetTenantVersionFromRequest(r)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusBadRequest)
		return
	}

	err = overview.CreateTemporaryOverview(userIdentity, version, subject)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
		return
	}

	session, err := database.OpenSession()
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
		return
	}
	defer session.Close()

	for {
		recordList, err := f(session)
		if err != nil {
			httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
			return
		}

		if len(recordList) == 0 {
			break // No more records to insert
		}

		err = overview.BulkInsertIntoOverview(userIdentity, subject, recordList, true)
		if err != nil {
			httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
			return
		}
	}

	err = overview.CommitOverview(userIdentity, subject)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
		return
	}

	httpcomm.ServiceResponse{
		Data: "OK",
	}.WriteData(w, httpcomm.PayloadFormatJSON)
}
