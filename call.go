package imosrpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func CallUrl(hostname string, method string, request interface{}, responsePtr interface{}) error {
	client := http.Client{}
	form, err := BuildForm(request)
	if err != nil {
		return fmt.Errorf("failed to build query: %s", err)
	}
	response, err := client.PostForm(hostname+"/"+method, form)
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
