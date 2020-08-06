package services

import (
	"Atrovan_Q1/db"
	"Atrovan_Q1/services/models"
	"Atrovan_Q1/validators"
	"net/http"
)

// CreateNewLesson service
// validate input
// create a lesson -> conflict checker
func CreateNewLesson(lesson *models.Lesson, userID int) (int, []byte) {
	var lessonModel db.LessonModel
	if ok, response := validators.Validate(lesson); !ok {
		return response.StatusCode, response.Body
	}
	if ok, response := lessonModel.InsertOne(lesson, userID); !ok {
		return response.StatusCode, response.Body
	}
	return http.StatusCreated, []byte("")
}
