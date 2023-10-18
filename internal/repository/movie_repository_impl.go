package repository

import (
	"database/sql"
	"discord-movie-service/models"
	"strings"
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

func (r *DatabaseMovieRepository) Get(params models.MovieQueryParams) ([]models.Movie, error) {

	var query strings.Builder
	query.WriteString("SELECT TitleType, Title, Genres, Year, Rating FROM movies WHERE 1=1") // The "1=1" is a dummy condition to simplify appending real conditions

	// Append conditions based on provided parameters
	if params.Genre != "" {
		query.WriteString(" AND genre = ?")
	}
	if params.Year != "" {
		query.WriteString(" AND year = ?")
	}
	if params.Rating != 0 {
		query.WriteString(" AND rating >= ?")
	}
	query.WriteString(" LIMIT 50") //TODO rethink this
	// Prepare the statement
	stmt, err := r.db.Prepare(query.String())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Collect arguments based on provided parameters
	var args []interface{}
	if params.Genre != "" {
		args = append(args, params.Genre)
	}
	if params.Year != "" {
		args = append(args, params.Year)
	}
	if params.Rating != 0 {
		args = append(args, params.Rating)
	}

	// Execute the query
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.TitleType, &movie.Title, &movie.Genres, &movie.Year, &movie.Rating); err != nil {
			return nil, err
		}
		movies = append(movies, movie)

		// Break if we've fetched the required amount of movies
		if params.Amount != 0 && len(movies) >= params.Amount {
			break
		}
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}
