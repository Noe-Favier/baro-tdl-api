package models

type CategoryUser struct {
	ID uint `json:"id" gorm:"primary_key"`

	CategoryID uint `json:"category"`
	UserID     uint `json:"user"`

	Category Category `json:"categoryId" gorm:"foreignKey:CategoryID;references:ID"`
	User     User     `json:"userId" gorm:"foreignKey:User;references:ID"`
}
