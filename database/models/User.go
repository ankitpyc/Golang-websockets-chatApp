package databases

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"userName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
