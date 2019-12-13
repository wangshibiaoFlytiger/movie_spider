package m_movie

import (
	"apiproject/entity"
	"github.com/globalsign/mgo/bson"
)

type Movie struct {
	Id            bson.ObjectId        `bson:"_id" json:"id"`
	DetailPageUrl string               `bson:"detailPageUrl" json:"detailPageUrl"`
	Title         string               `bson:"title" json:"title"`
	Actor         string               `bson:"actor" json:"actor"`
	CoverUrl      string               `bson:"coverUrl" json:"coverUrl"`
	Category      string               `bson:"category" json:"category"`
	Director      string               `bson:"director" json:"director"`
	PublishYear   int                  `bson:"publishYear" json:"publishYear"`
	Location      string               `bson:"location" json:"location"`
	Language      string               `bson:"language" json:"language"`
	Tag           string               `bson:"tag" json:"tag"`
	Desc          string               `bson:"desc" json:"desc"`
	PlayUrlInfo   map[string][]PlayUrl `bson:"playUrlInfo" json:"playUrlInfo"`
	//类型:1电影, 2电视剧
	Type int `bson:"type" json:"type"`

	CreateTime *entity.JsonTime `bson:"createTime" json:"createTime"`
	UpdateTime *entity.JsonTime `bson:"updateTime" json:"updateTime"`
	DeleteTime *entity.JsonTime `bson:"deleteTime" json:"deleteTime"`
}
