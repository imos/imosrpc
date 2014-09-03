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

func call(t *testing.T) {
	request := ExampleRequest{Value1: 100, Value2: 200}
	response := ExampleResponse{}
	if err := imosrpc.Call("call_test", request, &response); err != nil {
		t.Fatal(err)
	}
	expectedResponse := ExampleResponse{Addition: 300, Subtraction: -100}
	if !reflect.DeepEqual(expectedResponse, response) {
		t.Errorf("expected: %#v, actual: %#v.", expectedResponse, response)
	}
}

func TestCall(t *testing.T) {
	imosrpc.SetHostname(testServer.URL)
	call(t)
}

func TestInternalCall(t *testing.T) {
	imosrpc.SetHostname(imosrpc.InternalHostname)
	call(t)
}
