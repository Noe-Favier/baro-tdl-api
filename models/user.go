package models

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Username string `json:"username" gorm:"unique_index"`
	Email    string `json:"email"`
	Passwd   string `json:"passwd"`
	Roles    string `json:"roles"`
}
