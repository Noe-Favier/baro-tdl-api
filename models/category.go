package models

type Category struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Label string `json:"label"`
}
