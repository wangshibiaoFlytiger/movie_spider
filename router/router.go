package router

import (
	apiproject_config "apiproject/config"
	"apiproject/controller"
	c_kafka "apiproject/controller/kafka"
	c_movie "apiproject/controller/movie"
	c_video "apiproject/controller/video"
	c_wxpay "apiproject/controller/wxpay"
	"apiproject/middleware"
	rice "github.com/GeertJohan/go.rice"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"html/template"
)

/**
初始化路由
*/
func Init() *gin.Engine {
	engine := gin.Default()

	/**************************start 配置中间件 *****************/
	//支持跨域
	config := cors.DefaultConfig()
	config.AddAllowHeaders("X-Requested-With")
	config.AllowAllOrigins = true
	engine.Use(cors.New(config))
	//全局配置api日志访问中间件
	engine.Use(middleware.ApiLogMiddleware)
	/**************************end 配置中间件 *****************/

	/***********************start 通过go.rice配置页面模板 **********************/
	//配置模板文件的根目录
	templateBox := rice.MustFindBox("../public/template")

	//配置模板文件路径列表, 需填写相对于模板相对路径
	list := [...]string{"index.html"}
	for _, x := range list {
		templateString, err := templateBox.String(x)
		if err != nil {
			panic(err)
		}

		tmplMessage, err := template.New(x).Parse(templateString)
		if err != nil {
			panic(err)
		}

		engine.SetHTMLTemplate(tmplMessage)
	}
	/***********************end 配置页面模板 **********************/

	//通过go.rice配置静态文件目录
	engine.StaticFS("/static", rice.MustFindBox("../public/static").HTTPBox())

	//配置首页入口
	engine.GET("/", controller.Index)
	engine.GET("/index", controller.Index)

	//视频相关接口
	videoGroup := engine.Group("/api/video")
	//加密api响应数据中间件
	if apiproject_config.GlobalConfig.ServiceApiResponseEncrypt {
		videoGroup.Use(middleware.EncryptResponseMiddleware)
	}
	videoGroup.GET("/findList", c_video.FindVideoList)
	videoGroup.GET("/findVideoByWhere", c_video.FindVideoByWhere)
	videoGroup.GET("/addVideo", c_video.AddVideo)
	videoGroup.GET("/bulkAddVideo", c_video.BulkAddVideo)
	videoGroup.POST("/updateVideo", c_video.UpdateVideo)
	videoGroup.POST("/deleteVideo", c_video.DeleteVideo)

	//kafka相关
	kafkaGroup := engine.Group("/api/kafka")
	//加密api响应数据中间件
	if apiproject_config.GlobalConfig.ServiceApiResponseEncrypt {
		kafkaGroup.Use(middleware.EncryptResponseMiddleware)
	}
	kafkaGroup.GET("/sendMessage", c_kafka.SendMessage)

	//微信支付相关
	wxpayGroup := engine.Group("/api/wxpay")
	//加密api响应数据中间件
	if apiproject_config.GlobalConfig.ServiceApiResponseEncrypt {
		wxpayGroup.Use(middleware.EncryptResponseMiddleware)
	}
	wxpayGroup.POST("/wxH5Pay", c_wxpay.WxH5Pay)
	wxpayGroup.POST("/wxH5PayCallback", c_wxpay.WxH5PayCallback)

	spiderGroup := engine.Group("/api/spider")
	spiderGroup.POST("/craw88ysFilm", c_movie.Craw88ysFilm)

	movieGroup := engine.Group("/api/movie")
	movieGroup.GET("/findMovieList", c_movie.FindMovieList)
	movieGroup.GET("/getMovie", c_movie.GetMovie)

	return engine
}
