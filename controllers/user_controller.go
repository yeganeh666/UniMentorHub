package controllers

import (
	"Atrovan_Q1/db"
	"Atrovan_Q1/services"
	"Atrovan_Q1/services/models"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// Data is struct use for retrive lessons of user
type Data struct {
	Lessons []models.Lesson `json:"lessons"`
}

// CreateLesson creates a new lesson
// add lesson & make relation between lesson - user
func CreateLesson(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	lesson := models.Lesson{}
	if r.Body == nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&lesson); err != nil {
		errResponse, _ := json.Marshal(models.ErrorResponse{
			Code:      "BadArgument",
			Message:   "User input type is not valid",
			MessageFa: "جنس ورودی‌های کاربر معتبر نمی‌باشد.",
			Target:    []models.Target{{Name: "Fields", Description: err.Error()}},
		})
		response(w, http.StatusBadRequest, errResponse)
		return
	}

	userID, _ := strconv.Atoi(r.Header.Get("id"))
	responseStatus, res := services.CreateNewLesson(&lesson, userID)

	response(w, responseStatus, res)
}

// GetAllLessonsCache check chache returns all lessons.
func GetAllLessonsCache(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	page, limit := queryExtractor(r.URL.Query())
	from, to := page*limit, page*limit+limit
	redisConn, err := db.ConnectRedis()
	if err != nil {
		next(w, r)
	}
	data, err := redisConn.GetLessons(from, to)
	if err != nil {
		next(w, r)
		return
	}

	ok, lessons := unmarshalRedisData(data, to-from)
	if !ok {
		next(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(Data{Lessons: lessons})
	w.Write(response)
}

// GetAllLessons returns all lessons.
func GetAllLessons(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var lessonModel db.LessonModel
	page, limit := queryExtractor(r.URL.Query())

	data, err := lessonModel.FindAll(page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	response, _ := json.Marshal(Data{Lessons: *data})
	w.Write(response)
}

// GetUserAllLessons returns the lessons of specific user.
func GetUserAllLessons(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var userModel db.UserModel

	userID, _ := strconv.Atoi(r.Header.Get("id"))
	data, err := userModel.FindOne(&models.User{UserID: userID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	response, _ := json.Marshal(Data{Lessons: data.Lessons})
	w.Write(response)
}

// SignUpForLesson add lesson to user.
func SignUpForLesson(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var userModel db.UserModel
	lesson := models.Lesson{}
	if r.Body == nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&lesson); err != nil {
		errResponse, _ := json.Marshal(models.ErrorResponse{
			Code:      "BadArgument",
			Message:   "User input type is not valid",
			MessageFa: "جنس ورودی‌های کاربر معتبر نمی‌باشد.",
		})
		response(w, http.StatusBadRequest, errResponse)
		return
	}
	userID, _ := strconv.Atoi(r.Header.Get("id"))
	_, res := userModel.AddLessons(userID, lesson.Code)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.StatusCode)
	w.Write(res.Body)
}

func queryExtractor(query url.Values) (int, int) {
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		limit = 20
	}
	return page, limit
}

func unmarshalRedisData(data []string, count int) (bool, []models.Lesson) {
	lesson := models.Lesson{}
	lessons := []models.Lesson{}

	if len(data) == 0 || len(data) != count {
		return false, nil
	}

	for _, d := range data {
		if err := json.Unmarshal([]byte(d), &lesson); err != nil {
			return false, nil
		}
		lessons = append(lessons, lesson)
	}

	return true, lessons
}
