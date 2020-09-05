package models

type User struct {
	ID             int64  `json:"-,omitempty"`
	Username       string `json:"username,omitempty"`
	Email          string `json:"email,omitempty"`
	Role           string `json:"role,omitempty"`
	Password       string `json:"-"`
	HashedPassword string `json:"-"`
}
