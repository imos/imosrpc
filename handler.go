package imosrpc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
)

var rpcMethods = map[string]interface{}{}

func init() {
	http.HandleFunc("/", DefaultHandler)
}

func RegisterHandler(method string, handler interface{}) {
	method = strings.TrimPrefix(method, "/")
	handlerType := reflect.ValueOf(handler).Type()
	if handlerType.Kind() != reflect.Func {
		log.Fatalf("handler for method '%s' must be a function.", method)
	}
	if handlerType.IsVariadic() || handlerType.NumIn() != 1 ||
		handlerType.In(0).Kind() != reflect.Struct {
		log.Fatalf("method '%s' must take exactly 1 struct argument.", method)
	}
	if handlerType.NumOut() != 1 || handlerType.Out(0).Kind() != reflect.Struct {
		log.Fatalf("method '%s' must return exactly 1 struct value.", method)
	}
	rpcMethods[method] = handler
}

func getHandler(method string) (handler interface{}, ok bool) {
	method = strings.TrimPrefix(method, "/")
	handler, ok = rpcMethods[method]
	return
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	rpcHandler, ok := getHandler(r.URL.Path)
	w.Header().Set("Content-Type", "text/plain")
	if !ok {
		log.Printf("URL '%s' is not found.", r.URL.Path)
		w.WriteHeader(404)
		return
	}
	if err := r.ParseForm(); err != nil {
		w.Write([]byte(fmt.Sprintf(
			"failed to parse query as an HTTP request: %s", err)))
		w.WriteHeader(400)
		return
	}
	request := reflect.New(reflect.ValueOf(rpcHandler).Type().In(0))
	if err := ParseForm(request.Interface(), r.Form); err != nil {
		w.Write([]byte(fmt.Sprintf(
			"failed to parse query as a request: %s", err)))
		w.WriteHeader(400)
		return
	}
	response := reflect.ValueOf(rpcHandler).Call(
		[]reflect.Value{request.Elem()})[0].Interface()
	jsonData, err := json.Marshal(response)
	if err != nil {
		w.Write([]byte("failed to encode the response."))
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		w.WriteHeader(500)
	}
}
