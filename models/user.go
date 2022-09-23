package models

import "gorm.io/gorm"

type User struct {
	ID       uint   `gorm:"primary_key"`
	
	gorm.Model

	Username string `json:"username" gorm:"unique_index"`
	Email    string `json:"email"`
	Passwd   string `json:"passwd"`
	Roles    string `json:"roles" sql:"DEFAULT:'USER'"`

	//relations
	Categories []Category `gorm:"many2many:user_category;"`
}
