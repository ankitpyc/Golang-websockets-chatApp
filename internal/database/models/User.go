package databases

import (
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"userName"`
	Email    string `json:"email"`
	Password string `json:"password"`
	UserId   string `gorm:"-" json:"userid"`
}

// AfterCreate copies ID primary key tox userId
func (user *User) AfterCreate(tx *gorm.DB) (err error) {
	// Update the UserID field with the ID value
	user.UserId = fmt.Sprintf("%d", user.ID)
	// Save the updated user
	return nil
}

func (user *User) GetPassword() string {
	return user.Password
}

func (user *User) SetPassword(password string) {
	user.Password = password
}
