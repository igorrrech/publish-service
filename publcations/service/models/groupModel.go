package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Numder  uint
	Name    string
	Posts   []Post `gorm:"foreignKey:GroupID"`
	Parents []Group
}
