package forms

type FormCreateCategory struct {
	Label             string `json:"label"`
	CreatedByUsername string `json:"creator"`
}
