package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"turnstile/models"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

// // GetToken Api handler function, returns a
// func GetToken(db *gorm.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")

// 		fps, err := models.GetUsers(db)

// 		err = json.NewEncoder(w).Encode(fps)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 	}
// }

// CheckToken
func CheckToken(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		token := r.FormValue("token")

		tokenClaims, err := checkToken(token)

		if tokenClaims.Type != NormalTokenType {
			err = json.NewEncoder(w).Encode(&RespMessage{
				http.StatusBadRequest,
				"Bad token type"})
			return
		}

		if err != nil {
			err = json.NewEncoder(w).Encode(&RespMessage{
				http.StatusInternalServerError,
				err.Error()})
			return
		}

		// Get the time until the token expires
		diff := time.Until(time.Unix(tokenClaims.StandardClaims.ExpiresAt, 0))

		// result = time.Unix(claims.StandardClaims.ExpiresAt, 0).Format("2006-01-02 15:04:05")
		result := fmt.Sprintf("Expires in %02vm %02vs, Username %v", int(diff.Seconds())/60, int(diff.Seconds())%60, tokenClaims.UserName)

		err = json.NewEncoder(w).Encode(&RespMessage{
			http.StatusOK,
			fmt.Sprintf("%s", result)})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// RefreshToken
func RefreshToken(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		token := r.FormValue("token")

		tokenClaims, err := checkToken(token)

		if tokenClaims.Type != RefreshTokenType {
			err = json.NewEncoder(w).Encode(&LoginRespMessage{
				http.StatusBadRequest,
				"Bad token type",
				"",
				""})
			return
		}

		if err != nil {
			err = json.NewEncoder(w).Encode(&LoginRespMessage{
				http.StatusInternalServerError,
				err.Error(),
				"",
				""})
			return
		}

		usr, err := models.GetUserByUserName(db, tokenClaims.UserName)

		if err != nil {
			err = json.NewEncoder(w).Encode(&LoginRespMessage{
				http.StatusInternalServerError,
				err.Error(),
				"",
				""})
			return
		}

		token, refreshToken, err := CreateTokens(usr)

		// Send the token to the user
		err = json.NewEncoder(w).Encode(&LoginRespMessage{
			http.StatusOK,
			"Token refreshed successfully",
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

// CreateTokens
func CreateTokens(u models.User) (tokenString string, refreshTokenString string, err error) {

	// Create the Claims
	tokenClaims := MyCustomClaims{
		NormalTokenType,
		u.UserName,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer:    "turnstile",
		},
	}

	refreshTokenClaims := MyCustomClaims{
		RefreshTokenType,
		u.UserName,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
			Issuer:    "turnstile",
		},
	}

	// Create the token and refresh token and sign them
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	tokenString, err = token.SignedString(MySigningKey)

	// If one token fails, do not create the other
	if err != nil {
		return
	}

	refreshTokenString, err = refreshToken.SignedString(MySigningKey)

	return
}

func checkToken(tokenString string) (claims *MyCustomClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySigningKey, nil
	})

	claims, ok := token.Claims.(*MyCustomClaims)

	if !ok && !token.Valid {
		return
	}

	return
}
