package main

import (
	"discord-movie-service/internal/app"
	"discord-movie-service/internal/handlers"
	"discord-movie-service/internal/repository"
	"discord-movie-service/router"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//build table

	db := repository.NewDatabase()
	repository.SetupDatabaseSchema(db)

	//repo and handler init
	movieRepo := repository.NewDatabaseMovieRepository(db)
	movieHandler := handlers.NewMovieHandler(movieRepo)

	//initialize the app (this is just a list of the handlers)
	app := &app.App{MovieHandler: movieHandler}

	router := router.SetupRouter(app)
	router.Run(":8080")
	defer db.Close()
}
