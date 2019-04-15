package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
)

var config = struct {
	Type     string
	Host     string
	Port     uint16
	Name     string
	User     string
	Password string
	Path     string
}{
	"sqlite3",
	"localhost",
	3365,
	"app_database",
	"User",
	"Password",
	"data/turnstile.sqlite3",
}

// InitDatabase creates the database connection and creates the database
// structure if not created
func InitDatabase() (db *gorm.DB, err error) {
	switch config.Type {

	case "sqlite3":
		if _, err := os.Stat("data"); os.IsNotExist(err) {
			os.Mkdir("data", 0755)
		}
		db, err = gorm.Open(config.Type, config.Path)
	case "mysql":
		uri := "%s:%s@%s:%d/%s?charset=utf8&parseTime=True&loc=Local"
		db, err = gorm.Open(config.Type, fmt.Sprintf(uri, config.User,
			config.Password, config.Host, config.Port, config.Name))
	case "postgres":
		uri := "host=%s port=%d user=%s password=%s dbname=%s"
		db, err = gorm.Open(config.Type, fmt.Sprintf(uri, config.Host,
			config.Port, config.User, config.Password, config.Name))
	case "mssql":
		uri := "sqlserver://%s:%s@%s:%d?database=%s"
		db, err = gorm.Open(config.Type, fmt.Sprintf(uri, config.User,
			config.Password, config.Host, config.Port, config.Name))
	}

	if err != nil {
		errors.Wrap(err, "Cannot create the database")
		return
	}

	// User
	db.AutoMigrate(&User{})

	return
}
