package s_movie

import (
	d_movie "apiproject/dao/movie"
	"apiproject/log"
	m_movie "apiproject/model/movie"
	"apiproject/util"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/emirpasic/gods/utils"
	"github.com/globalsign/mgo/bson"
	"github.com/imroc/req"
	"go.uber.org/zap"
	"strings"
	"time"
)

var Movie88ysService = &movie88ysService{}

//爬取www.88ys.com网站的service
type movie88ysService struct {
}

/**
爬取www.88ys.com的影片
*/
func (this *movie88ysService) CrawFilmList() {
	for pageNo := 1; pageNo <= 902; pageNo++ {
		listPageUrl := "https://www.88ys.com/vod-type-id-1-pg-" + utils.ToString(pageNo) + ".html"
		//类型:1电影, 2电视剧
		movieType := 1
		this.CrawListPage(listPageUrl, movieType)
		log.Logger.Info("爬取www.88ys.com的影片", zap.Any("pageNo", pageNo), zap.Any("listPageUrl", listPageUrl))
	}

	log.Logger.Info("爬取www.88ys.com的影片, 完成")
}

/**
爬取影片列表
*/
func (this *movie88ysService) CrawListPage(listPageUrl string, movieType int) {
	document, err := goquery.NewDocument(listPageUrl)
	if err != nil {
		log.Logger.Error("爬取影片列表, 请求列表页异常", zap.Error(err))
		return
	}
	log.Logger.Info("爬取影片列表, 开始", zap.Any("listPageUrl", listPageUrl))

	//解析影片列表
	document.Find("div.index-area ul li").Each(func(i int, itemSelection *goquery.Selection) {
		log.Logger.Info("爬取影片列表, 解析影片列表, 开始", zap.Any("i", i))

		//解析列表页
		movie := m_movie.Movie{}
		movie.Type = movieType
		a := itemSelection.Find("a").First()
		href, exists := a.Attr("href")
		if !exists {
			println("测试爬取网站mov920, 异常1")
			log.Logger.Error("爬取影片列表, 解析列表页, href不存在")
			return
		}
		movie.DetailPageUrl = "https://www.88ys.com" + href

		lzbz := a.Find("span.lzbz").First()
		p1 := lzbz.Find("p").First()
		p2 := p1.Next()
		p3 := p2.Next()
		//p4 := p3.Next()

		movie.Title = p1.Text()
		movie.Category = p3.Text()

		//爬取详情页
		if err = this.CrawDetail(&movie); err != nil {
			log.Logger.Error("爬取影片列表, 爬取详情页异常", zap.Any("movie", movie), zap.Error(err))
			return
		}

		//入库
		if d_movie.MovieDao.ExistByCondition(bson.M{"title": movie.Title}) {
			log.Logger.Error("爬取影片列表, 已存在", zap.Any("title", movie.Title))
			return
		}
		if err = d_movie.MovieDao.Insert(&movie); err != nil {
			log.Logger.Error("爬取影片列表, 入库异常", zap.Any("movie", movie), zap.Error(err))
			return
		}

		log.Logger.Info("爬取影片列表, 解析影片列表, 完成", zap.Any("i", i), zap.Any("movie", movie))
		time.Sleep(2 * time.Second)
	})

	log.Logger.Info("爬取影片列表, 完成")
}

/**
爬取详情页
*/
func (this *movie88ysService) CrawDetail(movie *m_movie.Movie) error {
	document, err := goquery.NewDocument(movie.DetailPageUrl)
	if err != nil {
		log.Logger.Error("爬取详情页, 请求详情页url异常", zap.Any("url", movie.DetailPageUrl))
		return errors.New("请求详情页url异常")
	}

	leftDiv := document.Find("div.ct")
	cover := leftDiv.Find("div.ct-l img")
	coverUrl, exists := cover.Attr("src")
	if !exists {
		log.Logger.Error("爬取详情页, src不存在")
		return errors.New("src不存在")
	}
	movie.CoverUrl = coverUrl

	centerDiv := document.Find("div.ct-c")

	movie.Title = centerDiv.Find("dl h1").Text()
	dt1 := centerDiv.Find("dl dt").First()
	dt2 := dt1.Next()
	dd1 := centerDiv.Find("dl dd").First()
	dd2 := dd1.Next()
	dd3 := dd2.Next()
	dd4 := dd3.Next()
	dd5 := dd4.Next()
	dd6 := dd5.Next()

	movie.Actor = strings.ReplaceAll(dt2.Text(), "主演：", "")
	movie.Director = strings.ReplaceAll(dd3.Text(), "导演：", "")
	movie.PublishYear, _ = util.StrToInt(strings.ReplaceAll(dd5.Text(), "年份：", ""))
	movie.Location = strings.ReplaceAll(dd4.Text(), "地区：", "")
	movie.Language = strings.ReplaceAll(dd6.Text(), "语言：", "")
	movie.Tag = ""
	movie.Desc = strings.ReplaceAll(centerDiv.Find("div.ee").Text(), "剧情简介:                                　　", "")

	//查询播放源列表
	movie.PlayUrlInfo = make(map[string][]m_movie.PlayUrl)
	document.Find("div.playfrom li").Each(func(i int, selection *goquery.Selection) {
		sourceTitle := strings.TrimSpace(selection.Text())
		sourceId, _ := selection.Attr("id")

		//查询该数据源下的播放url列表
		playUrlList := []m_movie.PlayUrl{}
		playListDivSelection := "div.playlist[id='s" + sourceId + "'] a"
		document.Find(playListDivSelection).Each(func(i int, selection *goquery.Selection) {
			//访问播放页面
			playPageUrlTmp, exists := selection.Attr("href")
			if !exists {
				log.Logger.Error("爬取详情页, 解析播放页面url, href不存在")
				return
			}
			playPageUrlTmp = "https://www.88ys.com" + playPageUrlTmp
			html, err := util.GetDynamicPageHtmlContent(playPageUrlTmp)
			if err != nil {
				log.Logger.Error("爬取详情页, 解析播放页面url, 动态请求url异常", zap.Any("url", playPageUrlTmp), zap.Error(err))
				return
			}

			//解析播放页面url
			playPageDocument, err := goquery.NewDocumentFromReader(strings.NewReader(html))
			if err != nil {
				log.Logger.Error("爬取详情页, 解析播放页面url, NewDocumentFromReader异常", zap.Error(err))
				return
			}
			iframe1 := playPageDocument.Find("iframe").First()
			iframe2 := iframe1.Next()
			playPageUrl, exists := iframe2.Attr("src")
			if !exists {
				log.Logger.Error("爬取详情页, 解析播放页面url, src不存在", zap.Any("playPageUrlTmp", playPageUrlTmp))
				return
			}

			//解析视频url
			videoUrl := ""
			if strings.Contains(playPageUrl, "url=") {
				videoUrl = strings.Split(playPageUrl, "url=")[1]
			} else {
				resp, err := req.Get(playPageUrl, req.Header{"User-Agent": "Mozilla/5.0 (Linux; Android 8.0.0; LLD-AL10 Build/HONORLLD-AL10; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/63.0.3239.111 Mobile Safari/537.36;jxsafebrowser5620"})
				if err != nil {
					log.Logger.Error("爬取详情页, 解析播放页面url, 请求playPageUrl失败", zap.Error(err))
					return
				}
				split := strings.Split(resp.String(), "var main = \"")
				if len(split) < 2 {
					log.Logger.Error("爬取详情页, 解析播放页面url, 提取视频url异常", zap.Any("resp", resp.String()), zap.Any("playPageUrl", playPageUrl))
					return
				}
				tmp1 := split[1]
				videoUrl = "https://youku.com-ok-56.com" + strings.Split(tmp1, "\";")[0]
			}

			playUrl := m_movie.PlayUrl{
				Title:       selection.Text(),
				PlayPageUrl: playPageUrl,
				VideoUrl:    videoUrl,
			}
			playUrlList = append(playUrlList, playUrl)
		})

		movie.PlayUrlInfo[sourceTitle] = playUrlList
	})

	if len(movie.PlayUrlInfo) == 0 {
		return errors.New("playUrlList为空")
	}

	return nil
}
