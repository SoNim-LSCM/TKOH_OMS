package db_models

type Floors struct {
	FloorID    int     `json:"floorId"`
	FloorName  string  `json:"floorName"`
	FloorImage string  `json:"floorImage"`
	OriginX    float64 `json:"origin_x"`
	OriginY    float64 `json:"origin_y"`
	Resolution float64 `json:"resolution"`
	MapX       int     `json:"map_x"`
	MapY       int     `json:"map_y"`
}
