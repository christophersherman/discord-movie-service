package repository

import (
	"discord-movie-service/models"
)

type MovieRepository interface {
	Add(movie models.Movie) error
	Get(params models.MovieQueryParams) ([]models.Movie, error)
}
