package endpoint

import (
	"net/http"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/httpcomm"
	"github.com/gorilla/mux"
)

type onCreateCombobox func(userIdentity auth.UserIdentity, subject string) (any, error)

func GetComboboxBySubject(w http.ResponseWriter, r *http.Request, f onCreateCombobox) {
	vars := mux.Vars(r)
	subject := vars["subject"]

	if subject == "" {
		httpcomm.SetResponseError(&w, "no subject found in the request", nil, http.StatusBadRequest)
		return
	}

	userIdentity, err := auth.GetUserIdentityFromRequest(*r)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusUnauthorized)
		return
	}

	combobox, err := f(userIdentity, subject)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
	}

	httpcomm.ServiceResponse{
		Data: combobox,
	}.WriteData(w, httpcomm.PayloadFormatJSON)
}
