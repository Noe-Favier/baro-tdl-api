package models

import "gorm.io/gorm"

type Element struct {
	gorm.Model

	ID      uint   `json:"id" gorm:"primary_key"`
	Label   string `json:"label"`
	Checked bool   `json:"checked"`

	//relations
	CategoryID uint
}
