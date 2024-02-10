package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Client   *http.Client
	Request  *http.Request
	Response *http.Response
	Url      string
	Method   string
	Body     []byte
	//Body          *strings.Reader
	ResponseBody  []byte
	Timeout       time.Duration
	Authorization string
	Retry         int
	Debug         bool
}

type Request2 struct {
	Client   *http.Client
	Request  *http.Request
	Response *http.Response
	Url      string
	Method   string
	//Body     []byte
	Body *strings.Reader
	//Body          *strings.Reader
	ResponseBody  []byte
	Timeout       time.Duration
	Authorization string
	Retry         int
	Debug         bool
}

func (o *Request) Log(format string, v ...interface{}) {
	if o.Debug {
		log.Printf(format, v...)
	}
}

func (o *Request2) Log2(format string, v ...interface{}) {
	if o.Debug {
		log.Printf(format, v...)
	}
}

func (o *Request) Make() (err error) {
	//o.Request, err = http.NewRequest(o.Method, o.Url, bytes.NewBuffer(o.Body))
	//strings.NewReader(string(o.Body))
	o.Request, err = http.NewRequest(o.Method, o.Url, strings.NewReader(string(o.Body)))
	//o.Request, err = http.NewRequest(o.Method, o.Url, o.Body)
	if err != nil {
		o.Log("%v", err)
		return err
	}
	o.Request.Header.Set("Content-Type", "application/json")
	//o.Request.Header.Set("User-Agent", "terraform")
	if o.Authorization != "" {
		o.Request.Header.Set("Authorization", o.Authorization)
	}

	transport := &http.Transport{
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}
	o.Client = &http.Client{Transport: transport}
	o.Response, err = o.Client.Do(o.Request)
	if err != nil {
		o.Log("%v", err)
		return err
	}
	//var prettyJSON bytes.Buffer
	//if o.Method == "POST" || o.Method == "PATCH" || o.Method == "PUT" {
	//	err = json.Indent(&prettyJSON, o.Body, "", "\t")
	//	if err != nil {
	//		o.Log("JSON parse error: ", err)
	//		return err
	//	}
	//	o.Log("API Request: %s %d %s %s", o.Method, o.Response.StatusCode, o.Url, prettyJSON.String())
	//}
	o.Log("API Request: %s %d %s", o.Method, o.Response.StatusCode, o.Url)
	o.ResponseBody, err = io.ReadAll(o.Response.Body)
	_ = o.Response.Body.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	if (len(o.ResponseBody) == 0 || bytes.Equal(o.ResponseBody, []uint8{34, 34})) && (o.Method == "DELETE" || o.Method == "PATCH") {
		return nil
	}

	firstSymbol := string(o.ResponseBody[0:1])

	if firstSymbol != "{" && firstSymbol != "[" {
		return fmt.Errorf("API response not in json format:\n%s", o.ResponseBody)
	}

	var prettyJSONResp bytes.Buffer
	err = json.Indent(&prettyJSONResp, o.ResponseBody, "", "\t")
	if err != nil {
		o.Log("JSON parse error: ", err)
		return err
	}
	o.Log("API Response: %s", string(prettyJSONResp.Bytes()))
	return nil
}

func (o *Request2) Make2() (err error) {
	//o.Request, err = http.NewRequest(o.Method, o.Url, bytes.NewBuffer(o.Body))
	//strings.NewReader(string(o.Body))
	o.Request, err = http.NewRequest(o.Method, o.Url, o.Body)

	//o.Request, err = http.NewRequest(o.Method, o.Url, o.Body)
	if err != nil {
		o.Log2("%v", err)
		return err
	}

	//	//  -H 'accept: application/json' \
	//	//  -H 'Content-Type: application/x-www-form-urlencoded' \

	o.Request.Header.Set("accept", "application/json")
	//o.Request.Header.Set("Content-Type", "application/json")
	o.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//o.Request.Header.Set("User-Agent", "terraform")
	if o.Authorization != "" {
		o.Request.Header.Set("Authorization", o.Authorization)
	}

	transport := &http.Transport{
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}
	o.Client = &http.Client{Transport: transport}
	o.Response, err = o.Client.Do(o.Request)
	if err != nil {
		o.Log2("%v", err)
		return err
	}
	//var prettyJSON bytes.Buffer
	//if o.Method == "POST" || o.Method == "PATCH" || o.Method == "PUT" {
	//	err = json.Indent(&prettyJSON, o.Body, "", "\t")
	//	if err != nil {
	//		o.Log("JSON parse error: ", err)
	//		return err
	//	}
	//	o.Log("API Request: %s %d %s %s", o.Method, o.Response.StatusCode, o.Url, prettyJSON.String())
	//}
	o.Log2("API Request: %s %d %s", o.Method, o.Response.StatusCode, o.Url)
	o.ResponseBody, err = io.ReadAll(o.Response.Body)
	_ = o.Response.Body.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	if (len(o.ResponseBody) == 0 || bytes.Equal(o.ResponseBody, []uint8{34, 34})) && (o.Method == "DELETE" || o.Method == "PATCH") {
		return nil
	}

	firstSymbol := o.ResponseBody[0]
	if firstSymbol != 123 && firstSymbol != 91 {
		return fmt.Errorf("API response not in json format:\n%s", o.ResponseBody)
	}

	var prettyJSONResp bytes.Buffer
	err = json.Indent(&prettyJSONResp, o.ResponseBody, "", "\t")
	if err != nil {
		o.Log2("JSON parse error: ", err)
		return err
	}
	o.Log2("API Response: %s", string(prettyJSONResp.Bytes()))
	return nil
}
