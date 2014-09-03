package imosrpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Call(url string, request interface{}, responsePtr interface{}) error {
	client := http.Client{}
	form, err := BuildForm(request)
	if err != nil {
		return fmt.Errorf("failed to build query: %s", err)
	}
	response, err := client.PostForm(url, form)
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
