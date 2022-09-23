package models

import "gorm.io/gorm"

type User struct {
	ID uint `gorm:"primary_key" json:"-"`

	gorm.Model

	Username string `json:"username" gorm:"unique;not null"`
	Email    string `json:"email" gorm:"unique;not null"`

	Passwd string `json:"-"`
	Roles  string `json:"roles" sql:"DEFAULT:'USER'"`

	//relations
	Categories []Category `gorm:"many2many:user_category;"`
}
