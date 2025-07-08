package endpoint

import (
	"net/http"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/httpcomm"
	"github.com/gorilla/mux"
)

type onCreateCombobox func(userIdentity auth.UserIdentity, subject string, params map[string]string) (any, error)

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

	params := make(map[string]string)
	q := r.URL.Query()
	for k := range q {
		params[k] = q.Get(k)
	}

	combobox, err := f(userIdentity, subject, params)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
	}

	httpcomm.ServiceResponse{
		Data: combobox,
	}.WriteData(w, httpcomm.PayloadFormatJSON)
}
