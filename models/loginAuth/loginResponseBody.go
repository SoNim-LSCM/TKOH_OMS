package loginAuth

type LoginResponseBody struct {
	ID               int    `json:"userId"`
	Username         string `json:"username"`
	UserType         string `json:"userType"`
	AuthToken        string `json:"authToken"`
	LoginTime        string `json:"loginTime"`
	TokenExpiryTime  string `json:"tokenExpiryTime"`
	DutyLocationId   int    `json:"dutyLocationId"`
	DutyLocationName string `json:"dutyLocationName"`
}
