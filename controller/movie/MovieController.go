package c_movie

import (
	"apiproject/log"
	s_movie "apiproject/service/movie"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	type Para struct {
		Page int `form:"page" binding:"required"`
		Rows int `form:"rows" binding:"required"`
	}
	//绑定参数到对象
	para := Para{}
	if err := ctx.ShouldBind(&para); err != nil {
		log.Logger.Error("绑定请求参数到对象异常", zap.Error(err))
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": nil,
			"msg":  "参数错误",
		})
		return
	}
	log.Logger.Info("绑定请求参数到对象", zap.Any("channelNewsPara", para))

	movieList, totalCount := s_movie.MovieService.FindMovieList(para.Page, para.Rows)

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
