package controllers

import (
	"Atrovan_Q1/services"
	"Atrovan_Q1/services/models"
	"encoding/json"
	"net/http"
)

// CreateNewUser creates new user.
func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	if r.Body == nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		errResponse, _ := json.Marshal(models.ErrorResponse{
			Code:      "BadArgument",
			Message:   "User input type is not valid",
			MessageFa: "جنس ورودی‌های کاربر معتبر نمی‌باشد.",
			Target:    []models.Target{{Name: "Fields", Description: err.Error()}},
		})
		response(w, http.StatusBadRequest, errResponse)
		return
	}
	responseStatus, res := services.CreateNewUser(&user)
	response(w, responseStatus, res)
}

// Login handle user login to system.
func Login(w http.ResponseWriter, r *http.Request) {
	requestUser := models.User{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	responseStatus, token := services.Login(&requestUser)
	response(w, responseStatus, token)
}

// RefreshToken recreate the token
// notice : refresh token func not work after [exp] time is over, and user should login again in that state.
func RefreshToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	user := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&user)

	w.Header().Set("Content-Type", "application/json")
	w.Write(services.RefreshToken(user, r.Header.Get("role")))
}

// Logout handle logout process.
func Logout(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err := services.Logout(r)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
