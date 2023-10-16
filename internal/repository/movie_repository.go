package repository

import (
	"discord-movie-service/models"
)

type MovieRepository interface {
	Add(movie models.Movie) error
}
