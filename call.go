package imosrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const InternalHostname = "http://internal"

type responseWriter struct {
	ResponseData []byte
	StatusCode   int
}

func (r *responseWriter) Header() http.Header {
	return http.Header{}
}

func (r *responseWriter) Write(data []byte) (length int, err error) {
	length = len(data)
	err = nil
	r.ResponseData = append(r.ResponseData, data...)
	return
}

func (r *responseWriter) WriteHeader(code int) {
	r.StatusCode = code
}

func post(url string, data url.Values) (response *http.Response, err error) {
	client := http.Client{}
	if !strings.HasPrefix(url, InternalHostname+"/") {
		response, err = client.PostForm(url, data)
		return
	}
	request, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return
	}
	rw := &responseWriter{StatusCode: 200}
	DefaultHandler(rw, request)
	if rw.StatusCode != 200 {
		err = fmt.Errorf("status code is %d.", rw.StatusCode)
		return
	}
	response = &http.Response{
		Status:        fmt.Sprintf("%d", rw.StatusCode),
		StatusCode:    rw.StatusCode,
		Proto:         "HTTP",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewReader(rw.ResponseData)),
		ContentLength: int64(len(rw.ResponseData)),
		Close:         true,
		Request:       request,
	}
	return
}

func CallUrl(hostname string, method string, request interface{}, responsePtr interface{}) error {
	form, err := BuildForm(request)
	if err != nil {
		return fmt.Errorf("failed to build query: %s", err)
	}
	response, err := post(hostname+"/"+method, form)
	if err != nil {
		return fmt.Errorf("failed to request: %s", err)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read a response: %s", err)
	}
	err = json.Unmarshal(responseData, responsePtr)
	if err != nil {
		return fmt.Errorf("failed to parse a response: %s", err)
	}
	return nil
}

var defaultHostname string

func SetHostname(hostname string) error {
	defaultHostname = hostname
	return nil
}

func Call(method string, request interface{}, reseponsePtr interface{}) error {
	if defaultHostname == "" {
		return fmt.Errorf("hostname is not set.")
	}
	return CallUrl(defaultHostname, method, request, reseponsePtr)
}
