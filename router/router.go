package router

import (
	"discord-movie-service/internal/app"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router_app *app.App) *gin.Engine {
	r := gin.Default()

	//endpoints here

	r.POST("/addmovie", router_app.MovieHandler.AddMovie)
	r.GET("/movies", router_app.MovieHandler.GetMovie)
	return r
}
