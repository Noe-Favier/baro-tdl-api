package models

type Element struct {
	ID      uint   `json:"id" gorm:"primary_key"`
	Label   string `json:"label"`
	Checked bool   `json:"checked"`
}
