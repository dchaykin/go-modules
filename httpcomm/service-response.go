package httpcomm

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/dchaykin/go-modules/log"
)

type PayloadFormat string

const (
	PayloadFormatXML  PayloadFormat = "application/xml"
	PayloadFormatJSON PayloadFormat = "application/json"
)

func (pf PayloadFormat) IsXML() bool {
	return pf == PayloadFormatXML
}

func (pf PayloadFormat) IsJSON() bool {
	return !pf.IsXML()
}

func (pf PayloadFormat) String() string {
	return string(pf)
}

func GetAcceptHeader(req *http.Request) PayloadFormat {
	return PayloadFormat(req.Header.Get("Accept"))
}

func GetContentTypeHeader(req *http.Request) PayloadFormat {
	return PayloadFormat(req.Header.Get("Content-Type"))
}

type ServiceError struct {
	Code    *int   `json:"code,omitempty"`
	Message string `json:"message"`
}

type ServiceResponse struct {
	Data  interface{} `json:"data"`
	Error *string     `json:"error"`
}

func (sr *ServiceResponse) setError(msg string, err error) {
	sr.Error = new(string)
	if msg == "" && err != nil {
		*sr.Error = fmt.Sprintf("%v", err)
	} else if msg != "" && err == nil {
		*sr.Error = msg
	} else if msg != "" && err != nil {
		*sr.Error = fmt.Sprintf("%s: %v", msg, err)
	} else {
		sr.Error = nil
	}
}

func SetResponseError(w *http.ResponseWriter, msg string, err error, httpStatus int) {
	log.Error(msg, err)
	sr := ServiceResponse{}
	sr.setError(msg, err)
	data, e := json.Marshal(sr)
	if e != nil {
		http.Error(*w, fmt.Sprintf("%v", e), http.StatusInternalServerError)
	} else {
		http.Error(*w, string(data), httpStatus)
	}
}

func (resp ServiceResponse) WriteData(w http.ResponseWriter, format PayloadFormat) {
	b := new(bytes.Buffer)
	if format == PayloadFormatXML {
		e := xml.NewEncoder(b)
		e.Encode(resp)
		w.Header().Add("Content-Type", PayloadFormatXML.String())
	} else {
		e := json.NewEncoder(b)
		e.Encode(resp)
		w.Header().Add("Content-Type", PayloadFormatJSON.String())
	}
	fmt.Fprintf(w, "%v", b)
}
