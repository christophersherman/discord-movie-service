package app

import "discord-movie-service/internal/handlers"

// expand this with handlers as the project grows.
type App struct {
	MovieHandler *handlers.MovieHandler
}
