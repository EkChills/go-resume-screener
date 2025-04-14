package main

import (
	"log"

	"github.com/ekchills/go-resume-screener/database"
	"github.com/ekchills/go-resume-screener/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	err = database.ConnectDb()
	if err != nil {
		panic("failed to connect database")
	}
	database.MigrateDb()
	server := gin.Default()
	r := routes.Routes{Server: server}
	r.RegisterRoutes()
	err = server.Run(":8080")
	if err != nil {
		panic("failed to start server")
	}

}