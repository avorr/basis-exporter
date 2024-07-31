package client

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	ApiUrl       string
	clientId     string
	clientSecret string
	authUrl      string
)

func init() {
	ok := false
	clientId, ok = os.LookupEnv("BE_CLIENT_ID")
	if !ok || clientId == "" {
		panic("env var \"BE_CLIENT_ID\" isn't exist or = \"\"")
	}
	clientSecret, ok = os.LookupEnv("BE_CLIENT_SECRET")
	if !ok || clientSecret == "" {
		panic("env var \"BE_CLIENT_SECRET\" isn't exist or = \"\"")
	}
	authUrl, ok = os.LookupEnv("BE_AUTH_URL")
	if !ok || authUrl == "" {
		panic("env var \"BE_AUTH_URL\" isn't exist or = \"\"")
	}
	ApiUrl, ok = os.LookupEnv("BE_API_URL")
	if !ok || ApiUrl == "" {
		panic("env var \"BE_API_URL\" isn't exist or = \"\"")
	}
}

type Api struct {
	Host  string
	Token string
	Debug bool
}

func NewApi() *Api {
	debug := false
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		debug = true
	}
	return &Api{
		Host:  ApiUrl,
		Debug: debug,
	}
}

func (o *Api) NewRequest(method, resource string, body *strings.Reader, expectCode int) (*Request, error) {
	request := &Request{
		Debug:         o.Debug,
		Url:           fmt.Sprintf("%s/%s", o.Host, resource),
		Method:        method,
		Authorization: fmt.Sprintf("Bearer %s", o.Token),
	}

	if body != nil {
		request.Body = body
	}
	err := request.Make()
	if err != nil {
		return nil, err
	}
	if request.Response.StatusCode != expectCode {
		return request, fmt.Errorf(
			"wrong statusCode from API: %d, expect: %d, resource [%s], response: %s",
			request.Response.StatusCode,
			expectCode,
			resource,
			string(request.ResponseBody),
		)
	}
	return request, nil
}

func (o *Api) Auth() error {
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	paramString := fmt.Sprintf(
		"grant_type=client_credentials&client_id=%s&client_secret=%s&response_type=id_token",
		clientId,
		clientSecret,
	)
	var data = strings.NewReader(paramString)
	req, _ := http.NewRequest(http.MethodPost, authUrl, data)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Println(err)
		}
	}(resp.Body)
	token, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	o.Token = string(token)
	return nil
}

func (o *Api) NewRequestPost(url string, data *strings.Reader) ([]byte, int, error) {
	request, err := o.NewRequest("POST", url, data, 200)
	if err != nil {
		return nil, 0, err
	}
	return request.ResponseBody, request.Response.StatusCode, nil
}
