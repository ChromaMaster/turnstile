package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

var (
	// ErrUserDoesNotExist User does not exist
	errUserDoesNotExist = errors.New("User does not exist")
	// ErrUserAlreadyExists User already exists
	errUserAlreadyExists = errors.New("User already exists")
)

// User stores the user information
type User struct {
	gorm.Model

	UserName string `gorm:"UNIQUE"`
	Password string
}

// CreateUser stores the User in the database
func CreateUser(db *gorm.DB, u *User) (err error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	isExist := isUserExist(db, u.UserName)
	if isExist {
		return errUserAlreadyExists
	}

	if err = tx.Create(&u).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// GetUsers returns all the GetUsers stored in the database
func GetUsers(db *gorm.DB) ([]User, error) {
	var us []User
	if err := db.Find(&us).Error; err != nil {
		return us, errUserDoesNotExist
	}
	return us, nil
}

// GetUserID returns the User idenfified by the id
func GetUserByID(db *gorm.DB, id uint) (User, error) {
	var u User
	if err := db.Where(id).First(&u).Error; err != nil {
		return u, errUserDoesNotExist
	}
	return u, nil
}

// GetUserByUserName returns the User idenfified by the name
func GetUserByUserName(db *gorm.DB, userName string) (User, error) {
	var u User
	if err := db.Where("user_name = ?", userName).First(&u).Error; err != nil {
		return u, errUserDoesNotExist
	}
	return u, nil
}

func isUserExist(db *gorm.DB, userName string) bool {
	if err := db.Where("user_name = ?", userName).First(&User{}).Error; err != nil {
		return false
	}
	return true
}
