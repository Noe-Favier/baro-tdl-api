package forms

type FormCreateElement struct {
	Label             string `json:"label"`
	CreatedByUsername string `json:"creator"`
	CategoryCode      string `json:"code"`
}
