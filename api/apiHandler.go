package apiHandler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"tkoh_oms/config"
	errorHandler "tkoh_oms/errors"
)

var baseUrl string

func Init() {
	err := config.LoadENV()
	errorHandler.CheckError(err, "load env")
	baseUrl = os.Getenv("RFMS_URL")
}

func Get(endpoint string, param interface{}) ([]byte, error) {
	//Encode the data
	postBody, _ := json.Marshal(param)
	responseBody := bytes.NewBuffer(postBody)
	client := &http.Client{}
	body := []byte{}
	// response, err := http.Post(baseUrl+endpoint, "application/json", responseBody)
	req, err := http.NewRequest("GET", "http://"+baseUrl+endpoint, responseBody)
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth("oms", "YNetORtE")
	response, err := client.Do(req)
	if err != nil {
		return body, err
	}
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return body, err
	}

	return body, nil
}

func Post(endpoint string, param interface{}) ([]byte, error) {
	//Encode the data
	postBody, _ := json.Marshal(param)
	responseBody := bytes.NewBuffer(postBody)
	client := &http.Client{}
	body := []byte{}
	// response, err := http.Post(baseUrl+endpoint, "application/json", responseBody)
	req, err := http.NewRequest("POST", "http://"+baseUrl+endpoint, responseBody)
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth("oms", "YNetORtE")
	response, err := client.Do(req)
	if err != nil {
		return body, err
	}
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return body, err
	}

	return body, nil
}
