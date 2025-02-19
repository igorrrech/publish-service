package models

import (
	"errors"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	RawPost
}
type RawPost struct {
	GroupID uint   `json:"group-id"`
	UserID  uint   `json:"user-id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (p Post) Validate() error {
	if err := p.validateTitle(); err != nil {
		return err
	}
	if err := p.validateContent(); err != nil {
		return err
	}
	return nil
}
func (p Post) validateContent() error {
	return nil
}
func (p Post) validateTitle() error {
	return nil
}

var (
	ErrValidateContent = errors.New("content is not valide")
	ErrValidateTitle   = errors.New("title is not valise")
)
