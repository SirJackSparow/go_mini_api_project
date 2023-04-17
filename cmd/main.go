package main

import (
	"example/auth/api/controllers/authcontroller"
	"example/auth/mysql"
	"example/auth/usecase"
	"log"
	"net/http"

	"example/auth/repository"

	"example/auth/api/middlewares"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	mysql.ConnectDatabase()

	userRepo := repository.NewUserRepository(mysql.DB)
	userUseCase := usecase.NewUserUseCase(userRepo)
	authController := authcontroller.NewController(userUseCase)

	r.HandleFunc("/loginUser", authController.Login).Methods("POST")
	r.HandleFunc("/registerUser", authController.Register).Methods("POST")
	r.HandleFunc("/logout", authController.Logout).Methods("GET")

	route := r.PathPrefix("/route").Subrouter()
	route.HandleFunc("/profil", authController.ProfilUser).Methods("GET")
	route.Use(middlewares.JWTMiddleware)

	log.Fatal(http.ListenAndServe(":8080", r))

}
