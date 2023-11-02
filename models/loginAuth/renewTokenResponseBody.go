package loginAuth

type RenewTokenResponseBody struct {
	ID       int    `json:"userId"`
	Username string `json:"username" bson:"username"`
	// UserType            string `json:"userType" bson:"userType"`
	AuthToken           string `json:"authToken" bson:"authToken"`
	LoginDateTime       string `json:"loginDateTime" bson:"loginDateTime"`
	TokenExpiryDateTime string `json:"tokenExpiryDateTime" bson:"tokenExpiryDateTime"`
	// DutyLocationId      int    `json:"dutyLocationId" bson:"dutyLocationId"`
	// DutyLocationName    string `json:"dutyLocationName" bson:"dutyLocationName"`
}
