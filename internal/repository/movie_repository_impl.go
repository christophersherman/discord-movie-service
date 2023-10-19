package repository

import (
	"database/sql"
	"discord-movie-service/models"
	"math"
	"math/rand"
	"strings"
	"time"
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
	// Initialize the random number generator
	rand.Seed(time.Now().UnixNano())

	// Build query to fetch IDs based on filters
	var query strings.Builder
	query.WriteString("SELECT id FROM movies WHERE 1=1") // The "1=1" is a dummy condition to simplify appending real conditions

	if params.Genre != "" {
		query.WriteString(" AND MATCH(genres) AGAINST (? IN BOOLEAN MODE)")
	}
	if params.Year != "" {
		query.WriteString(" AND year  >= ?")
	}
	if params.Rating != 5.0 {
		query.WriteString(" AND rating >= ?")
	}

	if params.Runtime != 0 {
		if params.Runtime < -1 {
			query.WriteString(" AND runtime <= ?")
		} else {
			query.WriteString(" AND runtime >= ?")
		}
	}

	// Prepare the statement
	stmt, err := r.db.Prepare(query.String())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Collect arguments based on provided parameters
	var args []interface{}
	if params.Genre != "" {
		args = append(args, "+"+params.Genre+"*") // for LIKE query, we need to add % around the keyword
	}
	if params.Year != "" {
		args = append(args, params.Year)
	}
	if params.Rating != 5.0 {
		args = append(args, params.Rating)
	}
	if params.Runtime != 0 {
		args = append(args, int(math.Abs(float64(params.Runtime))))
	}

	// Execute the ID query
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect fetched IDs
	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return []models.Movie{}, nil // no movies found
	}

	// Randomly select params.Amount of IDs if there are more IDs than needed
	if len(ids) > params.Amount {
		rand.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })
		ids = ids[:params.Amount] // trim the slice to the needed amount
	}
	// Build query to fetch movies data based on selected IDs
	query.Reset()
	query.WriteString("SELECT TitleType, Title, Genres, Runtime, Year, Rating, Votes FROM movies WHERE id IN (")
	queryPlaceholders := strings.TrimRight(strings.Repeat("?,", len(ids)), ",") // creates ?,?,?... based on len(ids)
	query.WriteString(queryPlaceholders)
	query.WriteString(")")

	stmt, err = r.db.Prepare(query.String())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Convert ids to []interface{}
	args = make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	// Execute the query
	rows, err = stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.TitleType, &movie.Title, &movie.Genres, &movie.Runtime, &movie.Year, &movie.Rating, &movie.Votes); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	return movies, rows.Err()
}
