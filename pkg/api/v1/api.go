package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"

	"turnstile/pkg/api/v1/routes"

	"github.com/gorilla/mux"
)

// CreateAPIServer initializes the routes for the api
// It returns the router
func CreateAPIServer(db *gorm.DB) *mux.Router {
	router := mux.NewRouter()

	userRouter := router.PathPrefix("/user").Subrouter()
	// userRouter.HandleFunc("/", routes.GetUsers(db)).Methods("GET")
	// userRouter.HandleFunc("/{id}", routes.GetUser(db)).Methods("GET")

	userRouter.HandleFunc("/register", routes.RegisterUser(db)).Methods("POST")
	userRouter.HandleFunc("/login", routes.LoginUser(db)).Methods("POST")

	tokenRouter := router.PathPrefix("/token").Subrouter()
	tokenRouter.HandleFunc("/check", routes.CheckToken(db)).Methods("POST")
	tokenRouter.HandleFunc("/refresh", routes.RefreshToken(db)).Methods("POST")

	return router
}

// Serve will start the web server and will start polling
func Serve(router *mux.Router, address string, port uint32) {
	addressString := fmt.Sprintf("%s:%d", address, port)

	srv := &http.Server{
		Addr:         addressString,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      router,
	}

	log.Printf("Serving on %s", addressString)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
	defer srv.Close()
}
