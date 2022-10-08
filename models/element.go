package models

import "gorm.io/gorm"

type Element struct {
	gorm.Model

	ID      uint   `json:"-" gorm:"primary_key"`
	Label   string `json:"label"`
	Checked bool   `json:"checked"`
	Code    string `json:"code" gorm:"unique;not null"`

	CreatedByUsername string `json:"created_by_username"`

	//relations
	CategoryID uint
}
