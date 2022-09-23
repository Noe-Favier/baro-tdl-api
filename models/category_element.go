package models

type CategoryElement struct {
	ID uint `json:"id" gorm:"primary_key"`

	CategoryID uint
	ElementID  uint

	Category Category `json:"category" gorm:"foreignKey:ID;references:Category;constraint:"`
	Element  Element  `json:"element" gorm:"foreignKey:ID;references:Element;constraint:"`
}
