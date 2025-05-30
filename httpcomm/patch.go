package httpcomm

import (
	"io"
	"net/http"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/log"
)

func PatchToServiceIntern(serviceURL string, identity auth.UserIdentity, param map[string]string, payload []byte) (err error) {
	var response []byte
	if response, err = PatchData(serviceURL, identity, param, string(payload)); err != nil {
		if response != nil {
			log.Errorf("Response from %s: %s", serviceURL, string(response))
		}
		return err
	}

	serviceResponse, err := FetchServiceResponse(response)
	if err != nil {
		return err
	}

	if serviceResponse.Error != nil {
		log.Info("Response from %s: %v", serviceURL, *serviceResponse.Error)
	}

	return nil
}

func PatchData(endpoint string, identity auth.UserIdentity, parameters map[string]string, data ...string) (result []byte, err error) {
	hr := Patch(endpoint, identity, parameters, map[string]string{"Content-Type": getContentType()}, data...)
	return hr.Answer, hr.GetError()
}

func Patch(endpoint string, identity auth.UserIdentity, parameters map[string]string, headers map[string]string, data ...string) (httpResult HTTPResult) {
	payload := getPayloadFromSlice(data...)

	req, err := http.NewRequest("PATCH", endpoint, payload)
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
