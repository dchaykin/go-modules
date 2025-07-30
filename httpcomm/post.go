package httpcomm

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/log"
)

func Post(endpoint string, identity auth.UserIdentity, headers map[string]string, payload []byte) (httpResult HTTPResult) {
	if payload == nil {
		payload = []byte{}
	}
	log.Debug("/POST %s [ %s ]", endpoint, string(payload))
	return post(endpoint, false, identity, headers, bytes.NewBuffer(payload))
}

func PostBuffer(endpoint string, identity auth.UserIdentity, headers map[string]string, body *bytes.Buffer) (httpResult HTTPResult) {
	return post(endpoint, false, identity, headers, body)
}

func post(endpoint string, insecure bool, identity auth.UserIdentity, headers map[string]string, body *bytes.Buffer) (httpResult HTTPResult) {
	req, err := http.NewRequest("POST", endpoint, body)
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

func getPayloadFromSlice(data ...string) string {
	var payload string
	for _, d := range data {
		payload += d + " "
	}
	return payload
}
