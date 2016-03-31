package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type MibigClient struct {
	Host string
}

type MibigResponse struct {
	Error       bool   `json:"error"`
	Message     string `json:"message"`
	RedirectUrl string `json:"redirect_url"`
}

func (mc *MibigClient) ServiceInfo() (string, error) {
	uri := mc.Host + "/v1.0.0/bgc-registration"

	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("ServiceInfo got unexpected StatusCode of " + resp.Status)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func (mc *MibigClient) StoreMibigSubmissionOnly(raw string, version int) (*MibigResponse, error) {
	uri := mc.Host + "/v1.0.0/bgc-registration"

	resp, err := http.PostForm(uri,
		url.Values{"json": {raw}, "version": {fmt.Sprintf("%d", version)}})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r MibigResponse

	if err := json.Unmarshal(bodyBytes, &r); err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return &r, errors.New("StoreMibigSubmissionOnly got unexpected StatusCode of " + resp.Status)
	}

	return &r, nil
}

var RedirectError = errors.New("Don't follow HTTP redirect")

func (mc *MibigClient) StoreMibigSubmission(raw string, version int) (int, error) {
	uri := mc.Host + "/v2.0.0/bgc-registration"

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return RedirectError
		},
	}

	resp, err := client.PostForm(uri,
		url.Values{"json": {raw}, "version": {fmt.Sprintf("%d", version)}})
	if url_err, ok := err.(*url.Error); ok && url_err.Err == RedirectError {
		err = nil
	}
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 303 {
		return resp.StatusCode, errors.New("StoreMibigSubmission got unexpected StatusCode of " + resp.Status)
	}

	return resp.StatusCode, nil
}
