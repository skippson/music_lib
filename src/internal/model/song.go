package model

type Song struct {
	ID    uint   `gorm:"primary_key"`
	Group_name string `json:"group_name"`
	Song  string `json:"song"`
	// ReleaseDate string
	// Text        string
}
