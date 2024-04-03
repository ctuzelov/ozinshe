package models

import "time"

type User struct {
	ID            int       `json:"id,omitempty"`
	Name          string    `json:"name,omitempty" validate:"required,min=2,max=50"`
	Email         string    `json:"email,omitempty" validate:"required,email"`
	Number        string    `json:"number,omitempty" validate:"required"`
	DateOfBirth   time.Time `json:"date_of_birth,omitempty" validate:"required"`
	UserType      string    `json:"user_type,omitempty" validate:"required"`
	Password      string    `json:"password,omitempty" validate:"required,min=6"`
	Token         string    `json:"token,omitempty"`
	Refresh_Token string    `json:"refresh_token,omitempty"`
	User_id       string    `json:"user_id,omitempty"`
}
