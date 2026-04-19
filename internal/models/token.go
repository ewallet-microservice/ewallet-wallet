package models

type TokenData struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Fullname string `json:"full_name"`
	Email    string `json:"email"`
}
