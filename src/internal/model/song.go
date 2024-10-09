package model

type Song struct {
	ID          uint   `gorm:"primary_key"`
	Group_name  string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Lyrics      string `json:"text"`
}
