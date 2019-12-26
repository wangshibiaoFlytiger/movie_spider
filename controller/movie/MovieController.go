package c_movie

import (
	s_movie "apiproject/service/movie"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
爬取www.mov920.com站点的film
*/
func Craw88ysFilm(ctx *gin.Context) {
	go s_movie.Movie88ysService.CrawFilmList()

	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": nil,
	})
}

/**
查询影片列表
*/
func FindMovieList(ctx *gin.Context) {
	movieList, totalCount := s_movie.MovieService.FindMovieList()

	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": gin.H{
			"list":       movieList,
			"totalCount": totalCount,
		},
	})
}

/**
查询某影片
*/
func GetMovie(ctx *gin.Context) {
	movieId := ctx.Query("movieId")
	movie := s_movie.MovieService.GetMovie(movieId)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": movie,
	})
}
