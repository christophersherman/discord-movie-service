package handlers

import (
	"bytes"
	"discord-movie-service/internal/repository"
	"discord-movie-service/models"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	repo repository.MovieRepository
}

func NewMovieHandler(r repository.MovieRepository) *MovieHandler {
	return &MovieHandler{repo: r}
}

func (m *MovieHandler) AddMovie(c *gin.Context) {
	var movie models.Movie

	// Read the raw request body
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error reading raw body:", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	// Important: You need to reset the body if you want to read it again (for example, for binding)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	if err := c.BindJSON(&movie); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Println("Binding error:", err)
		log.Println("Raw JSON payload received:", string(rawBody))
		return
	}

	//swap to movie repo implementation
	err = m.repo.Add(movie)

	if err != nil && strings.Contains(err.Error(), "1062") || strings.Contains(err.Error(), "Duplicate entry") {
		// Handle the duplicate error
		c.JSON(409, gin.H{"error": "The movie already exists."}) // 409 Conflict might be a more suitable status code
	} else {
		// Handle other errors
		c.JSON(500, gin.H{"error": err.Error()})
		log.Println("Error:", err.Error())
	}

	c.JSON(200, gin.H{"message": "Movie added successfully!"})
}

func (m *MovieHandler) GetMovie(c *gin.Context) { //TODO clean this ugly function up

	movie_amount_str := c.DefaultQuery("amount", "1")
	movie_amount, err := strconv.Atoi(movie_amount_str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount parameter; it must be an integer."})
		return
	}

	genre_str := c.DefaultQuery("genre", "")
	year := c.DefaultQuery("year", "")
	runtime := c.DefaultQuery("runtime", "0")
	rating_str := c.DefaultQuery("rating", "5.0")
	rating, err := strconv.ParseFloat(rating_str, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating parameter; must be a float."})
		return
	}

	var movieQuery models.MovieQueryParams
	movieQuery.Amount = movie_amount
	movieQuery.Genre = genre_str
	movieQuery.Year = year
	movieQuery.Rating = float32(rating)

	movieQuery.Runtime, err = strconv.Atoi(runtime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid runtime parameter, must be an integer"})
	}

	movies, err := m.repo.Get(movieQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"movies": movies})
}
