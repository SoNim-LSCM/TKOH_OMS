package loginAuth

type LoginResponseBody struct {
	ID                  int    `json:"userId"`
	Username            string `json:"username"`
	UserType            string `json:"userType"`
	AuthToken           string `json:"authToken"`
	LoginDateTime       string `json:"loginDateTime"`
	TokenExpiryDateTime string `json:"tokenExpiryDateTime"`
	DutyLocationId      int    `json:"dutyLocationId"`
	DutyLocationName    string `json:"dutyLocationName"`
}
