package httpcomm

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	neturl "net/url"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/log"
)

type HTTPResult struct {
	Answer     []byte
	StatusCode int
	Status     string
	url        string
	method     string
	err        error
}

func (hr *HTTPResult) GetError() error {
	if hr.err != nil {
		return hr.err
	}
	switch hr.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		{
			return nil
		}
	}
	return fmt.Errorf("%s %s: %s (%d)", hr.method, hr.GetURL(), hr.Status, hr.StatusCode)
}

func (hr *HTTPResult) GetURL() string {
	u, err := neturl.Parse(hr.url)
	if err != nil {
		log.WrapError(err)
		return ""
	}
	return fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)
}

func GetData(url string) ([]byte, error) {
	hr := Get(url, nil, nil)
	return hr.Answer, hr.GetError()
}

func GetInsecure(url string, parameters map[string]string, headers map[string]string) (httpResult HTTPResult) {
	return get(url, parameters, headers, true)
}

func Get(url string, parameters map[string]string, headers map[string]string) (httpResult HTTPResult) {
	return get(url, parameters, headers, false)
}

func get(url string, parameters map[string]string, headers map[string]string, insecure bool) (httpResult HTTPResult) {
	log.Debug("/GET %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.WrapError(err)
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
	if insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	response, err := client.Do(req)
	if err != nil {
		return HTTPResult{err: err}
	}
	defer response.Body.Close()

	hr := HTTPResult{
		StatusCode: response.StatusCode,
		Status:     response.Status,
		url:        url,
		method:     "GET",
	}

	if hr.GetError() == nil {
		if hr.Answer, err = io.ReadAll(response.Body); err != nil {
			return HTTPResult{err: err}
		}
	}

	return hr
}

func GetFromServiceIntern(identity auth.UserIdentity, serviceURL string, param map[string]string) (serviceResponse ServiceResponse, err error) {
	var response []byte
	if response, err = getDataIntern(identity, serviceURL, param); err != nil {
		if response != nil {
			log.Errorf("Response from %s: %s", serviceURL, string(response))
		}
		serviceResponse.setError("", err)
		return serviceResponse, err
	}

	if serviceResponse, err = FetchServiceResponse(response); err != nil {
		return serviceResponse, err
	}

	if serviceResponse.Error != nil {
		log.Info("Response from %s: %v", serviceURL, *serviceResponse.Error)
	}

	return serviceResponse, nil
}

func getContentType() string {
	return "application/json"
}

func getDataIntern(identity auth.UserIdentity, endpoint string, parameters map[string]string) (result []byte, err error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err = identity.Set(req); err != nil {
		return nil, err
	}
	req.Header.Set("Accept", getContentType())

	q := req.URL.Query()
	for key, value := range parameters {
		q.Add(key, value)
	}

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if result, err = io.ReadAll(resp.Body); err != nil {
		return nil, err
	}

	return result, nil
}
