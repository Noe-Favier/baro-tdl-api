package forms

type FormLoginUser struct {
	Login    string `json:"login"` //can be username | email
	Password string `json:"password"`
}
