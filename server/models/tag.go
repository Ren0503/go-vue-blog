package models

import (
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type Tag struct {
	gorm.Model
	Name        string `gorm:"unique_index"`
	Slug        string `gorm:"unique_index"`
	Description string
	Articles    []Article `gorm:"many2many:articles_tags;"`
	IsNewRecord bool      `gorm:"-;default:false"`
}

func (a *Tag) BeforeSave() (err error) {
	a.Slug = slug.Make(a.Name)
	return
}
