package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"turnstile/models"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"golang.org/x/crypto/bcrypt"
)

// GetUsers
func GetUsers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		fps, err := models.GetUsers(db)

		err = json.NewEncoder(w).Encode(fps)
		if err != nil {
			err = json.NewEncoder(w).Encode(&RespMessage{
				http.StatusInternalServerError,
				err.Error()})
		}
	}
}

// GetUser
func GetUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)

		id, _ := strconv.ParseUint(params["id"], 10, 64)
		myFp, err := models.GetUserByID(db, uint(id))
		if err != nil {
			err = json.NewEncoder(w).Encode(&RespMessage{
				http.StatusBadRequest,
				err.Error()})
			return
		}

		err = json.NewEncoder(w).Encode(myFp)
		if err != nil {
			err = json.NewEncoder(w).Encode(&RespMessage{
				http.StatusInternalServerError,
				err.Error()})
		}
	}
}

// RegisterUser
func RegisterUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		username := r.FormValue("username")
		password := r.FormValue("password")

		hashedPassword, err := hashPassword(password)

		// Error hashing the password
		if err != nil {
			err = json.NewEncoder(w).Encode(&RespMessage{
				http.StatusInternalServerError,
				err.Error()})
			return
		}

		err = models.CreateUser(db, &models.User{
			UserName: username,
			Password: hashedPassword,
		})

		// Error creating the user
		if err != nil {
			err = json.NewEncoder(w).Encode(&RespMessage{
				http.StatusInternalServerError,
				err.Error()})
			return
		}

		err = json.NewEncoder(w).Encode(&RespMessage{
			http.StatusOK,
			"Successfully registered"})

		if err != nil {
			err = json.NewEncoder(w).Encode(&RespMessage{
				http.StatusInternalServerError,
				err.Error()})
		}
	}
}

// LoginUser
func LoginUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		username := r.FormValue("username")
		password := r.FormValue("password")

		usr, err := models.GetUserByUserName(db, username)

		if err != nil {
			err = json.NewEncoder(w).Encode(&LoginRespMessage{
				http.StatusUnauthorized,
				err.Error(),
				"",
				""})
			return
		}

		// Check password
		authorized := checkPasswordHash(password, usr.Password)

		if !authorized {
			err = json.NewEncoder(w).Encode(&LoginRespMessage{
				http.StatusUnauthorized,
				"Unable to login",
				"",
				""})
		}

		token, refreshToken, err := CreateTokens(usr)

		if err != nil {
			err = json.NewEncoder(w).Encode(&LoginRespMessage{
				http.StatusInternalServerError,
				err.Error(),
				"",
				""})
		}

		// Send the token to the user
		err = json.NewEncoder(w).Encode(&LoginRespMessage{
			http.StatusOK,
			"Logged in successfully",
			token,
			refreshToken})

		if err != nil {
			err = json.NewEncoder(w).Encode(&LoginRespMessage{
				http.StatusInternalServerError,
				err.Error(),
				"",
				""})
		}
	}
}

func hashPassword(password string) (hashedPassword string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	hashedPassword = string(bytes)
	return
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
