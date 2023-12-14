package rfms

type GetLocationResponse struct {
	ResponseCode    int    `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Body            struct {
		LocationList []struct {
			LocationID   int       `json:"locationId"`
			LocationName string    `json:"locationName"`
			Position     []float64 `json:"position"`
			Orientation  []float64 `json:"orientation"`
		} `json:"locationList"`
	} `json:"body"`
}
