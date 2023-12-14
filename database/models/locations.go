package db_models

type Locations struct {
	ID           int    `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	LocationID   int    `json:"locationId"`
	LocationName string `json:"locationName"`
}
