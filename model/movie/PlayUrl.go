package m_movie

type PlayUrl struct {
	Title       string `bson:"title" json:"title"`
	PlayPageUrl string `bson:"playPageUrl" json:"playPageUrl"`
	VideoUrl    string `bson:"videoUrl" json:"videoUrl"`
}
