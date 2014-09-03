package imosrpc_test

import (
	"github.com/imos/imosrpc"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var testServer *httptest.Server = nil

func init() {
	imosrpc.RegisterHandler("call_test", ExampleHandler)
	testServer = httptest.NewServer(http.HandlerFunc(imosrpc.DefaultHandler))
}

func TestCall(t *testing.T) {
	request := ExampleRequest{Value1: 100, Value2: 200}
	response := ExampleResponse{}
	if err := imosrpc.Call(testServer.URL+"/call_test", request, &response); err != nil {
		t.Fatal(err)
	}
	expectedResponse := ExampleResponse{Addition: 300, Subtraction: -100}
	if !reflect.DeepEqual(expectedResponse, response) {
		t.Errorf("expected: %#v, actual: %#v.", expectedResponse, response)
	}
}
