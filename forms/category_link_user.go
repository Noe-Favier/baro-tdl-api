package forms

type FormLinkCategoryToUser struct {
	Usernames    []string `json:"usernames"` //leave empty array to remove all users (except creator)
	CategoryCode string   `json:"category_code"`
}
