package ws_model

type SubscribeTokenResponse struct {
	MessageCode string `json:"messageCode"`
	UserID      int    `json:"userId"`
	Username    string `json:"username"`
	UserType    string `json:"userType"`
}
