package routes

import (
	"crypto-users/pkg/handlers"
	"crypto-users/pkg/models"
	"net/http"
)

func CreateRoutes(router *models.Router) (*models.Router) {
	router.HandleFunc("/", handlers.Home).Methods("GET")

	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	router.HandleFunc("/user/register", handlers.Register).Methods("POST")
	router.HandleFunc("/user/login", handlers.Login).Methods("POST")
	router.HandleFunc("/user/logout", handlers.Logout).Methods("GET")

	router.HandleFunc("/token/verify", handlers.Verify).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(handlers.NotFound)

	return router
}
