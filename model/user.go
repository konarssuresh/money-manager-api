package model

type User struct {
	UserId    string `json:"user_id"`
	Password  string `json:"password"`
	CreatedOn string `json:"created_on"`
	LastLogin string `json:"last_login"`
}
