package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model

	ID    uint   `json:"id" gorm:"primary_key"`
	Label string `json:"label"`
	Code  string `json:"code" gorm:"unique;not null"`

	CreatedByUsername string `json:"created_by_username"`

	//relations
	Elements []Element
	Users    []User `gorm:"many2many:user_category;"`
}
