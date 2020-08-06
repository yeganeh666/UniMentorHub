package db

import (
	"Atrovan_Q1/config"
	"Atrovan_Q1/services/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"

	// postgres driver for gorm.
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// LessonModel is a type to intract with db on lessons table.
type LessonModel struct {
}

// UserModel is a type to intract with db on users table.
type UserModel struct {
}

// connextGorm initialize the connection of orm to postgres.
func connectGorm() (*GormCli, error) {
	if gormInstance == nil {
		gormInstance = &GormCli{}
		username := config.EnvGetStr("POSTGRES_USER", "postgres")
		password := config.EnvGetStr("POSTGRES_PASSWORD", "postgres")
		dbName := config.EnvGetStr("POSTGRES_DB", "atrovan")
		dbHost := config.EnvGetStr("POSTGRES_DATABASE_HOST", "localhost")

		dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Build connection string

		conn, err := gorm.Open("postgres", dbURI)
		if err != nil {
			return nil, err
		}
		conn.AutoMigrate(&models.User{}, &models.Lesson{}, &models.Date{})
		gormInstance.conn = conn
		return gormInstance, nil
	}
	return gormInstance, nil
}

// InsertOne insert one user to database
// check existence -> unique_violation (default postgres lib)
func (userModel UserModel) InsertOne(user *models.User) (bool, *models.Response) {
	db, _ := connectGorm()
	err := db.conn.Create(&user).Error
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			response, _ := json.Marshal(models.ErrorResponse{Code: "Conflict", Message: "user already exists.", MessageFa: "کاربر با مشخصات وارد شده وجود دارد."})
			return false, &models.Response{StatusCode: http.StatusConflict, Body: response}
		}
		return false, &models.ServerDownResponse
	}

	if ok := db.conn.NewRecord(user); ok != false {
		return false, &models.ServerDownResponse
	}
	return true, &models.Response{}
}

// FindOne return users based on userId.
func (userModel UserModel) FindOne(user *models.User) (*models.User, error) {
	db, _ := connectGorm()
	var returnUser models.User

	if db.conn.Preload("Lessons", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Date").Find(&returnUser.Lessons)
	}).First(&returnUser, "user_id = ?", user.UserID).RecordNotFound() {
		return nil, fmt.Errorf("")
	}

	return &returnUser, nil
}

// FindAllLessons return user all lessons.
func (userModel UserModel) FindAllLessons(UserID int) (*models.User, error) {
	db, _ := connectGorm()
	var returnUser models.User
	err := db.conn.Preload("Lessons").First(&returnUser, "user_id = ?", UserID).Error
	if err != nil {
		return nil, err
	}
	return &returnUser, nil
}

// AddLessons add lessons to user.
func (userModel UserModel) AddLessons(userID int, id int) (bool, *models.Response) {
	var lessonModel LessonModel
	db, _ := connectGorm()

	user, err := userModel.FindOne(&models.User{UserID: userID})
	if err != nil {
		response, _ := json.Marshal(models.ErrorResponse{Code: "NotFound", Message: "user not found.", MessageFa: "کاربری با مشخصات وارد شده وجود ندارد."})
		return false, &models.Response{StatusCode: http.StatusNotFound, Body: response}
	}

	lesson, err := lessonModel.FindOne(id)
	if err != nil {
		response, _ := json.Marshal(models.ErrorResponse{Code: "NotFound", Message: "one or more than one of lessons is not find", MessageFa: "یک یا چند تا از کد درس‌هایی که وارد نموده‌اید پیدا نشدند"})
		return false, &models.Response{StatusCode: http.StatusNotFound, Body: response}
	}

	if lesson.Capacity == 0 {
		response, _ := json.Marshal(models.ErrorResponse{Code: "Full", Message: "class is full.", MessageFa: "ظرفیت کلاس مورد نظر پر شده است"})
		return false, &models.Response{StatusCode: http.StatusBadRequest, Body: response}

	}

	for _, les := range user.Lessons {
		for _, date := range les.Date {
			if ok, response := isTimeConflict(date, lesson); ok {
				return false, &models.Response{StatusCode: http.StatusNotFound, Body: response}
			}
		}
	}
	queryCreateRelation := fmt.Sprintf("INSERT INTO user_lessons(user_user_id, lesson_code) VALUES (%v, %v)", user.UserID, id)
	queryDecrementFromCapacity := fmt.Sprintf("UPDATE lessons SET capacity = capacity - 1 WHERE code %v", id)
	db.conn.Exec(queryCreateRelation)
	db.conn.Exec(queryDecrementFromCapacity)
	updateRedisCache(id)
	return true, &models.Response{StatusCode: http.StatusOK, Body: []byte("")}
}

// InsertOne insert one lesson to database.
func (lessonModel LessonModel) InsertOne(lesson *models.Lesson, userID int) (bool, *models.Response) {
	var userModel UserModel
	db, _ := connectGorm()
	user, err := userModel.FindOne(&models.User{UserID: userID})
	if err != nil {
		response, _ := json.Marshal(models.ErrorResponse{Code: "NotFound", Message: "user not found.", MessageFa: "کاربری با مشخصات وارد شده وجود ندارد."})
		return false, &models.Response{StatusCode: http.StatusNotFound, Body: response}
	}
	for _, less := range user.Lessons {
		for _, date := range less.Date {
			if ok, response := isTimeConflict(date, lesson); ok {
				return false, &models.Response{StatusCode: http.StatusBadRequest, Body: response}
			}
		}
	}
	if err = db.conn.Model(&user).Association("Lessons").Append(lesson).Error; err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			response, _ := json.Marshal(models.ErrorResponse{Code: "Conflict", Message: "lesson already exists.", MessageFa: "درس با مشخصات وارد شده وجود دارد."})
			return false, &models.Response{StatusCode: http.StatusConflict, Body: response}
		}
		return false, &models.ServerDownResponse
	}

	return true, &models.Response{}
}

// FindAll find all lessons from database.
func (lessonModel LessonModel) FindAll(page, limit int) (*[]models.Lesson, error) {
	db, _ := connectGorm()
	var returnLessons []models.Lesson
	if err := db.conn.Offset(page * limit).Limit(limit).Preload("Date").Find(&returnLessons).Error; err != nil {
		return nil, err
	}

	redisConn, _ := ConnectRedis()
	redisConn.AddLesson(returnLessons)

	return &returnLessons, nil
}

//FindOne find one lesson from database by given code.
func (lessonModel LessonModel) FindOne(id int) (*models.Lesson, error) {
	db, _ := connectGorm()
	var returnLesson models.Lesson
	if err := db.conn.Find(&returnLesson, "code = ?", id).Error; err != nil {
		return nil, err
	}

	return &returnLesson, nil
}

func isTimeConflict(d models.Date, lesson *models.Lesson) (bool, []byte) {
	timeStartM, _ := time.Parse(time.Kitchen, d.Start)
	timeEndM, _ := time.Parse(time.Kitchen, d.End)

	if timeStartM.After(timeEndM) {
		response, _ := json.Marshal(models.ErrorResponse{Code: "Time", Message: "start time should before end time.", MessageFa: "زمان شروع باید قبل از زمان پایان باشد."})
		return true, response
	}

	for _, dateL := range lesson.Date {
		if dateL.Day == d.Day {
			timeStartL, _ := time.Parse(time.Kitchen, dateL.Start)
			timeEndL, _ := time.Parse(time.Kitchen, dateL.End)
			if timeEndM.After(timeStartL) && timeEndL.After(timeStartM) {
				response, _ := json.Marshal(models.ErrorResponse{Code: "Conflict", Message: "time conflict.", MessageFa: "تداخل زمانی، شما قبلا در این زمان درس ارائه داده‌اید"})
				return true, response
			}
		}
	}
	return false, nil
}

func updateRedisCache(id int) {
	redis, _ := ConnectRedis()
	redis.conn.ZRem(ctx, lessonsKey, id)
}
