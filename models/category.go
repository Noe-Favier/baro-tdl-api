package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model

	ID    uint   `json:"-" gorm:"primary_key"`
	Label string `json:"label"`
	Code  string `json:"code" gorm:"unique;not null"`

	CreatedByUsername string `json:"creator"`

	//relations
	Elements []Element `json:"elements"`
	Users    []User    `json:"users" gorm:"many2many:user_category;"`
}
