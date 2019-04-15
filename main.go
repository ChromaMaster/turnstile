package main

import (
	"turnstile/models"
	"turnstile/pkg/api/v1"

	"github.com/jinzhu/gorm"
)

func main() {

	db, err := models.InitDatabase()

	if err != nil {
		panic(err)
	}

	populateDatabase(db)

	router := api.CreateAPIServer(db)
	api.Serve(router, "localhost", 8000)
}

func populateDatabase(db *gorm.DB) {
	u1 := models.User{
		UserName: "D3stru7",
		Password: "1234",
	}
	u2 := models.User{
		UserName: "ChromaMaster",
		Password: "6789",
	}
	u3 := models.User{
		UserName: "Tembleking",
		Password: "10111213",
	}

	models.CreateUser(db, &u1)
	models.CreateUser(db, &u2)
	models.CreateUser(db, &u3)

}
