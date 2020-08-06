package models

import "github.com/jinzhu/gorm"

type Date struct {
	LessonID uint   `json:"-" gorm:"AUTO_INCREMENT;primary_key"`
	Start    string `json:"start" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	End      string `json:"end" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Day      int    `json:"day" validate:"required,min=0,max=6"`
}

// Lesson define model for lessons.
type Lesson struct {
	gorm.Model
	Name     string `json:"name" validate:"required"`
	Code     int    `json:"code" gorm:"primary_key:true;unique_index;not null" validate:"required"`
	Unit     int    `json:"unit" validate:"required"`
	Date     []Date `json:"dates" validate:"required" gorm:"many2many:lesson_dates;association_foreignkey:lesson_id;foreignkey:code"`
	Capacity int    `json:"capacity" validate:"required"`
}

// TableName return lessons table name.
func (lesson *Lesson) TableName() string {
	return "lessons"
}

// TableName return lessons table name.
func (data *Date) TableName() string {
	return "dates"
}
