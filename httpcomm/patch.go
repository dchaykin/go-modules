package httpcomm

import (
	"bytes"
	"io"
	"net/http"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/log"
)

func Patch(endpoint string, identity auth.UserIdentity, parameters map[string]string, headers map[string]string, data ...string) (httpResult HTTPResult) {
	payload := getPayloadFromSlice(data...)

	log.Debug("/PATCH %s [ %s ]", endpoint, payload)

	req, err := http.NewRequest("PATCH", endpoint, bytes.NewReader([]byte(payload)))
	if err != nil {
		return HTTPResult{err: err}
	}

	if identity != nil {
		if err = identity.Set(req); err != nil {
			return HTTPResult{err: err}
		}
	}

	for key := range headers {
		req.Header.Set(key, headers[key])
	}

	q := req.URL.Query()
	for key, value := range parameters {
		q.Add(key, value)
	}

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return HTTPResult{err: err}
	}
	defer resp.Body.Close()

	hr := HTTPResult{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		url:        endpoint,
		method:     "PATCH",
	}

	if hr.GetError() == nil {
		if hr.Answer, err = io.ReadAll(resp.Body); err != nil {
			return HTTPResult{err: err}
		}
	}

	return hr
}
