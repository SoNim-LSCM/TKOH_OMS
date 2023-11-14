package loginAuth

type RenewTokenResponseBody struct {
	ID       int    `json:"userId"`
	Username string `json:"username" bson:"username"`
	// UserType            string `json:"userType" bson:"userType"`
	AuthToken       string `json:"authToken" bson:"authToken"`
	LoginTime       string `json:"loginTime" bson:"loginTime"`
	TokenExpiryTime string `json:"tokenExpiryTime" bson:"tokenExpiryTime"`
	// DutyLocationId      int    `json:"dutyLocationId" bson:"dutyLocationId"`
	// DutyLocationName    string `json:"dutyLocationName" bson:"dutyLocationName"`
}
