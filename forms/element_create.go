package forms

type FormCreateElement struct {
	Label             string `json:"label"`
	CreatedByUsername string `json:"created_by_username"`
	CategoryCode      string `json:"code"`
}
