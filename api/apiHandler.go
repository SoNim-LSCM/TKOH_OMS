package apiHandler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/SoNim-LSCM/TKOH_OMS/config"
	errorHandler "github.com/SoNim-LSCM/TKOH_OMS/errors"
)

var baseUrl string

func Init() {
	err := config.LoadENV()
	errorHandler.CheckError(err, "load env")
	baseUrl = os.Getenv("RFMS_URL")
}

func Get(endpoint string, param interface{}) []byte {
	response, err := http.Get(baseUrl + endpoint)
	if err != nil {
		log.Println(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	return body
}

func POST(endpoint string, param interface{}) []byte {
	//Encode the data
	postBody, _ := json.Marshal(param)
	responseBody := bytes.NewBuffer(postBody)
	client := &http.Client{}
	// response, err := http.Post(baseUrl+endpoint, "application/json", responseBody)
	req, err := http.NewRequest("POST", "http://"+baseUrl+endpoint, responseBody)
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth("oms", "YNetORtE")
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	return body
}
