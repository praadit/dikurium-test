// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

type CreateTodoInput struct {
	Title string `json:"title"`
}

type SigninInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninResult struct {
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

type SignupInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResult struct {
	Email string `json:"email"`
}
