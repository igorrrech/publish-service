package models

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Phone    string `gorm:"unique;index"`
	Password string
	Role     string
}

const (
	RoleAdmin   = "admin"
	RoleManager = "manager"
	RoleUser    = "user"
)

func (u *User) Validate() error {
	if err := u.validatePhone(); err != nil {
		return err
	}
	if err := u.validateRole(); err != nil {
		return err
	}
	return nil
}
func (u *User) validatePhone() error {
	return nil
}
func (u *User) validateRole() error {
	return nil
}

var (
	ErrValidatePhone = fmt.Errorf("phone is not valide")
	ErrValidateRole  = fmt.Errorf("role is not valide")
)
