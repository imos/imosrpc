package imosrpc_test

import (
	"github.com/imos/imosrpc"
	"net"
	"reflect"
	"testing"
)

var httpListener net.Listener

func init() {
	imosrpc.RegisterHandler("server_test", ExampleHandler)
	var err error
	httpListener, err = net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	imosrpc.SetHttpListener(httpListener)
	go func() {
		imosrpc.Serve()
	}()
}

func TestServe(t *testing.T) {
	request := ExampleRequest{Value1: 100, Value2: 200}
	response := ExampleResponse{}
	if err := imosrpc.CallUrl(
		"http://"+httpListener.Addr().String(), "server_test",
		request, &response); err != nil {
		t.Fatal(err)
	}
	expectedResponse := ExampleResponse{Addition: 300, Subtraction: -100}
	if !reflect.DeepEqual(expectedResponse, response) {
		t.Errorf("expected: %#v, actual: %#v.", expectedResponse, response)
	}
}
