package imosrpc_test

import (
	"encoding/json"
	"github.com/imos/imosrpc"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type ExampleRequest struct {
	Value1 int `json:"value1"`
	Value2 int `json:"value2"`
}

type ExampleResponse struct {
	Addition    int `json:"addition"`
	Subtraction int `json:"subtraction"`
}

func ExampleHandler(request ExampleRequest) ExampleResponse {
	return ExampleResponse{
		Addition:    request.Value1 + request.Value2,
		Subtraction: request.Value1 - request.Value2,
	}
}

type HttpResponse struct {
	ResponseData []byte
	ResponseCode int
}

func (r *HttpResponse) Header() http.Header {
	return http.Header{}
}

func (r *HttpResponse) Write(data []byte) (length int, err error) {
	length = len(data)
	err = nil
	r.ResponseData = append(r.ResponseData, data...)
	return
}

func (r *HttpResponse) WriteHeader(code int) {
	r.ResponseCode = code
}

func init() {
	imosrpc.RegisterHandler("example", ExampleHandler)
}

func TestPostMethod(t *testing.T) {
	request, err := http.NewRequest(
		"POST", "http://example.com/example",
		ioutil.NopCloser(strings.NewReader("value1=100&value2=200")))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatalf("failed to create a request: %s", err)
	}
	response := HttpResponse{}
	imosrpc.DefaultHandler(&response, request)
	expectedResponse := ExampleResponse{Addition: 300, Subtraction: -100}
	actualResponse := ExampleResponse{}
	json.Unmarshal(response.ResponseData, &actualResponse)
	if !reflect.DeepEqual(expectedResponse, actualResponse) {
		t.Errorf("%#v is expected, but %#v.", expectedResponse, actualResponse)
	}
}

func TestGetMethod(t *testing.T) {
	request, err := http.NewRequest(
		"GET", "http://example.com/example?value1=100&value2=200", nil)
	if err != nil {
		t.Fatalf("failed to create a request: %s", err)
	}
	response := HttpResponse{}
	imosrpc.DefaultHandler(&response, request)
	expectedResponse := ExampleResponse{Addition: 300, Subtraction: -100}
	actualResponse := ExampleResponse{}
	json.Unmarshal(response.ResponseData, &actualResponse)
	if !reflect.DeepEqual(expectedResponse, actualResponse) {
		t.Errorf("%#v is expected, but %#v.", expectedResponse, actualResponse)
	}
}

func BenchmarkGet(b *testing.B) {
	request, err := http.NewRequest(
		"GET", "http://api/example?value1=100&value2=200", nil)
	if err != nil {
		b.Fatalf("failed to create a request: %s", err)
	}
	for i := 0; i < b.N; i++ {
		response := httptest.NewRecorder()
		imosrpc.DefaultHandler(response, request)
	}
}
