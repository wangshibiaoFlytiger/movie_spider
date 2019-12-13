package d_movie

import (
	"apiproject/config"
	"apiproject/dao"
	"apiproject/entity"
	"apiproject/log"
	m_movie "apiproject/model/movie"
	"apiproject/mongo"
	"github.com/globalsign/mgo/bson"
	"go.uber.org/zap"
	"time"
)

var MovieDao = &movieDao{}

type movieDao struct {
	dao.BaseDao
}

/**
插入movie对象
*/
func (this *movieDao) Insert(movie *m_movie.Movie) error {
	sessionClone := mongo.GetMongoSessionClone()
	defer sessionClone.Close()
	db := sessionClone.DB(config.GlobalConfig.MongoDatabase)

	collection := db.C("movie")
	movie.Id = bson.NewObjectId()
	movie.CreateTime = &entity.JsonTime{time.Now()}
	movie.UpdateTime = movie.CreateTime
	err := collection.Insert(movie)
	if err != nil {
		log.Logger.Error("插入movie对象, 异常", zap.Any("movie", movie), zap.Error(err))
		return err
	}

	log.Logger.Info("插入movie对象, 完成", zap.Any("movie", movie))
	return nil
}

/**
查询movie列表
*/
func (this *movieDao) FindMovieList(condition bson.M, page int, rows int) (movieList []m_movie.Movie, totalCount int) {
	sessionClone := mongo.GetMongoSessionClone()
	defer sessionClone.Close()
	db := sessionClone.DB(config.GlobalConfig.MongoDatabase)

	collection := db.C("movie")
	err := collection.Find(condition).Sort("-publishYear").Skip((page - 1) * rows).Limit(rows).All(&movieList)
	if err != nil {
		log.Logger.Error("查询movie列表, 异常", zap.Any("condition", condition), zap.Error(err))
		return nil, 0
	}

	totalCount, err = collection.Find(condition).Count()
	if err != nil {
		log.Logger.Error("查询movie列表, 查询总数异常", zap.Error(err))
		return nil, 0
	}

	log.Logger.Info("查询movie列表, 完成", zap.Any("condition", condition), zap.Any("count", len(movieList)))
	return movieList, totalCount
}

/**
查询某movie
*/
func (this *movieDao) GetMovie(movieId string) (movie m_movie.Movie) {
	sessionClone := mongo.GetMongoSessionClone()
	defer sessionClone.Close()
	db := sessionClone.DB(config.GlobalConfig.MongoDatabase)

	collection := db.C("movie")
	err := collection.Find(bson.M{"_id": bson.ObjectIdHex(movieId)}).One(&movie)
	if err != nil {
		log.Logger.Error("查询某movie, 异常", zap.Any("movieId", movieId), zap.Error(err))
	}

	return movie
}

/**
按条件查询movie是否存在
*/
func (this *movieDao) ExistByCondition(condition bson.M) (exist bool) {
	sessionClone := mongo.GetMongoSessionClone()
	defer sessionClone.Close()
	db := sessionClone.DB(config.GlobalConfig.MongoDatabase)

	collection := db.C("movie")

	count, err := collection.Find(condition).Count()
	if err != nil {
		log.Logger.Error("按条件查询movie是否存在, 异常", zap.Any("condition", condition), zap.Error(err))
		return false
	}

	return count > 0
}

/**
根据条件查询某movie
*/
func (this *movieDao) GetMovieByCondition(condition bson.M) (movie m_movie.Movie) {
	sessionClone := mongo.GetMongoSessionClone()
	defer sessionClone.Close()
	db := sessionClone.DB(config.GlobalConfig.MongoDatabase)

	collection := db.C("movie")
	err := collection.Find(condition).One(&movie)
	if err != nil {
		log.Logger.Error("根据条件查询某movie, 异常", zap.Any("condition", condition), zap.Error(err))
	}

	return movie
}
