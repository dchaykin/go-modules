package httpcomm

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/log"
)

func PostData(endpoint string, identity auth.UserIdentity, data ...string) (result []byte, err error) {
	hr := post(endpoint, false, identity, map[string]string{"Content-Type": getContentType()}, data...)
	return hr.Answer, hr.GetError()
}

func PostInsecure(endpoint string, identity auth.UserIdentity, headers map[string]string, data ...string) (httpResult HTTPResult) {
	return post(endpoint, true, identity, headers, data...)
}

func Post(endpoint string, identity auth.UserIdentity, headers map[string]string, data ...string) (httpResult HTTPResult) {
	return post(endpoint, false, identity, headers, data...)
}

func post(endpoint string, insecure bool, identity auth.UserIdentity, headers map[string]string, data ...string) (httpResult HTTPResult) {
	payload := getPayloadFromSlice(data...)

	req, err := http.NewRequest("POST", endpoint, payload)
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
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	if insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return HTTPResult{err: err}
	}
	defer resp.Body.Close()

	hr := HTTPResult{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		url:        endpoint,
		method:     "POST",
	}

	if hr.GetError() == nil {
		if hr.Answer, err = io.ReadAll(resp.Body); err != nil {
			return HTTPResult{err: err}
		}
	}

	return hr
}

func PostToServiceIntern(serviceURL string, identity auth.UserIdentity, payload []byte) (err error) {
	var response []byte
	if response, err = PostData(serviceURL, identity, string(payload)); err != nil {
		if response != nil {
			log.Error("Response from %s: %s", serviceURL, string(response))
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

func getPayloadFromSlice(data ...string) *bytes.Reader {
	var payload string
	for _, d := range data {
		payload += d + " "
	}
	return bytes.NewReader([]byte(payload))
}
