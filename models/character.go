package models

import "time"

type (
	User struct {
		ID       int64  `json:"user_id" db:"user_id"`
		Email    string `json:"email" db:"email"`
		Fullname string `json:"fullname" db:"fullname"`
		Referral string `json:"referral" db:"referral"`
		Password string `json:"password" db:"password"`
		Token    string `json:"token"`
	}
	GetCharactersParam struct {
		Name           string    `json:"name"`
		NameStartsWith string    `json:"name_starts_with"`
		Comics         string    `json:"comics"`
		Series         string    `json:"series"`
		Events         string    `json:"events"`
		Stories        string    `json:"stories"`
		Limit          int64     `json:"limit"`
		Offset         int64     `json:"offset"`
		OrderBy        string    `json:"order_by"`
		ModifiedSince  time.Time `json:"modified_since"`
	}
)
