package models

import "github.com/jinzhu/gorm"

// User define model for students and professors
type User struct {
	gorm.Model
	FirstName string   `json:"first_name" validate:"required"`
	LastName  string   `json:"last_name" validate:"required"`
	Role      string   `json:"role" validate:"required,eq=professor|eq=student"`
	UserID    int      `json:"userId" gorm:"primary_key:true;unique_index;not null" validate:"required"`
	Password  string   `json:"password" validate:"required"`
	Lessons   []Lesson `json:"lessons" gorm:"many2many:user_lessons;association_foreignkey:code;foreignkey:user_id"`
}

// TableName return user table name.
func (user *User) TableName() string {
	return "users"
}
