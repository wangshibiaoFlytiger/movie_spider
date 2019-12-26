package s_movie

import (
	d_movie "apiproject/dao/movie"
	m_movie "apiproject/model/movie"
	"github.com/globalsign/mgo/bson"
)

var MovieService = &movieService{}

//爬取www.88ys.com网站的service
type movieService struct {
}

/**
查询影片列表
*/
func (this *movieService) FindMovieList() (movieList []m_movie.Movie, totalCount int) {
	movieList, totalCount = d_movie.MovieDao.FindMovieList(bson.M{}, 1, 1000)
	return
}

/**
查询某影片
*/
func (this *movieService) GetMovie(movieId string) (movie m_movie.Movie) {
	movie = d_movie.MovieDao.GetMovie(movieId)
	return movie
}
