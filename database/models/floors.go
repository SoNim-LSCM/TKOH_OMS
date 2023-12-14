package db_models

type Floors struct {
	FloorID    int    `json:"floorId"`
	FloorName  string `json:"floorName"`
	FloorImage string `json:"floorImage"`
}
