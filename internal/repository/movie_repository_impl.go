package repository

import (
	"database/sql"
	"discord-movie-service/models"
)

type DatabaseMovieRepository struct {
	db *sql.DB
}

func NewDatabaseMovieRepository(db *sql.DB) *DatabaseMovieRepository {
	return &DatabaseMovieRepository{db}
}

func (r *DatabaseMovieRepository) Add(movie models.Movie) error {

	_, err := r.db.Exec("INSERT INTO movies (titleType, title, year, genres, runtime, rating, votes) VALUES (?, ?, ?, ?, ?, ?, ?)", movie.TitleType, movie.Title, movie.Year, movie.Genres, movie.Runtime, movie.Rating, movie.Votes)

	if err != nil {
		return err
	}

	return nil
}
