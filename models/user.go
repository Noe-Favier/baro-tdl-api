package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	ID       uint   `json:"id" gorm:"primary_key"`
	Username string `json:"username" gorm:"unique_index"`
	Email    string `json:"email"`
	Passwd   string `json:"passwd"`
	Roles    string `json:"roles"`

	//relations
	Categories []Category `gorm:"many2many:user_category;"`
}
