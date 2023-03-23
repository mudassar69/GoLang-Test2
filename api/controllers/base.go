package controllers

import (
	"assignment2/api/middlewares"
	"assignment2/api/models"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Initialize() {
	var err error
	//connStr := "user=mudassarraza dbname=testdb sslmode=disabled"
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disabled", os.Getenv("DB_USER"), os.Getenv("DB_NAME"))
	a.DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		fmt.Println("Cannot connect to database")
	}

	a.DB.Debug().AutoMigrate(&models.User{})
	a.DB.Debug().AutoMigrate(&models.Task{})

	a.Router = mux.NewRouter()
	a.InitializeRoutes()
}

func (a *App) InitializeRoutes() {

	a.Router.Use(middlewares.SetContentTypeMiddleware)
	a.Router.HandleFunc("/register", a.RegisterUser).Methods("POST")
	a.Router.HandleFunc("/login", a.LoginUser).Methods("POST")

	subRouter := a.Router.PathPrefix("/api").Subrouter()
	subRouter.Use(middlewares.AuthJwtVerify)

	subRouter.HandleFunc("/task", a.AddTask).Methods("POST")
	subRouter.HandleFunc("/task", a.GetTasks).Methods("GET")
	subRouter.HandleFunc("/task/{id}", a.DeleteTask).Methods("DELETE")
	subRouter.HandleFunc("/task/{id}", a.UpdateTask).Methods("PUT")
	subRouter.HandleFunc("/task/assign", a.AssignTaskToUser).Methods("POST")
}

func (a *App) RunServer() {
	http.ListenAndServe(":5000", a.Router)
}
