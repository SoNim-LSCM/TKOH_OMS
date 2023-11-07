package dto

type AddDeliveryOrderDTO struct {
	OrderType             string `json:"orderType"`
	NumberOfAmrRequire    int    `json:"numberOfAmrRequire"`
	StartLocationID       int    `json:"startLocationId"`
	StartLocationName     string `json:"startLocationName"`
	ExpectingStartTime    string `json:"expectingStartTime"`
	EndLocationID         int    `json:"endLocationId"`
	EndLocationName       string `json:"endLocationName"`
	ExpectingDeliveryTime string `json:"expectingDeliveryTime"`
}
