package crawler

import (
	"strings"
	"io/ioutil"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"encoding/json"
	"net/url"
	"net/http"
	"log"
	"github.com/psy-core/MVCrawler/entity"
)

// 获取种子列表
func getSeedUrls() []string {
	area := []string{"ML", "US", "KR"}
	artist := []string{"Boy", "Girl"}
	version := []string{"music_video", "official_video", "subtitle", "fan_video", "fan_make_video", "yinyuetai"}
	tag := []string{"HyperCrystal", "HDV", "VchartOnTime", "HotCover", "ClassicMv", "SH"}
	genre := []string{"Dance", "Electronic", "Soundtrack", "Pop"}

	seedURLs := make([]string, 0)
	for i := range area {
		for j := range artist {
			for k := range version {
				for l := range tag {
					for m := range genre {
						seedURLs = append(seedURLs, "http://mv.yinyuetai.com/all?pageType=page&sort=pubdate&area="+area[i]+"&artist="+artist[j]+"&version="+version[k]+"&tag="+tag[l]+"&genre="+genre[m])
					}
				}
			}
		}
	}
	return seedURLs
}

// 根据种子列表获取mv info请求url列表
func getMvInfoUrls(limit int, seedURLs ...string) []string {
	mvJsonUrls := make([]string, 0)
	mvIdMap := make(map[string]bool)
	count := 0
	for i := range seedURLs {
		doc, err := goquery.NewDocument(seedURLs[i])
		if err != nil {
			log.Printf("url %s download and parse failed. ignore.", seedURLs[i])
			continue
		}

		// Find the review items
		doc.Find("#mvlist li").First().Each(func(j int, s *goquery.Selection) {
			// For each item found, get the band and title
			mvUrl, _ := s.Find("div.info>p>a").Attr("href")
			id := mvUrl[strings.LastIndex(mvUrl, "/")+1:]
			if id != "" && !mvIdMap[id] {
				count ++
				log.Printf("find mv id : %s ...\n", id)
				mvIdMap[id] = true
				mvJsonUrls = append(mvJsonUrls, "http://www.yinyuetai.com/insite/get-video-info?json=true&videoId="+id)
			}
		})

		if limit != -1 && count >= limit {
			break
		}
	}
	return mvJsonUrls
}

//获取MV对象列表
func getMVs(mvJsonUrls ...string) []entity.Mv {

	mvs := make([]entity.Mv, 0)
	for _, mvUrl := range mvJsonUrls {
		response, err := http.Get(mvUrl)
		if err != nil {
			log.Println("error! http get failed.")
			panic(err)
		}
		jsonbyte, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println("error! http response body read failed.")
			panic(err)
		}
		mvInfo := entity.MvInfo{}
		json.Unmarshal(jsonbyte, &mvInfo)

		mv := entity.Mv{
			Artist:   mvInfo.VideoInfo.CoreVideoInfo.ArtistName,
			Name:     mvInfo.VideoInfo.CoreVideoInfo.VideoName,
			AudioUrl: "",
		}

		videoUrlModels := mvInfo.VideoInfo.CoreVideoInfo.VideoUrlModels
		if len(videoUrlModels) <= 0 {
			continue
		}
		//log.Printf("videoUrlModels length: %d\n", len(videoUrlModels))
		for i := len(videoUrlModels) - 1; i >= 0; i-- {
			if videoUrlModels[i].QualityLevel == "sh" || videoUrlModels[i].QualityLevel == "he" {
				mv.AudioUrl = videoUrlModels[i].VideoUrl[:strings.LastIndex(videoUrlModels[i].VideoUrl, "?")]
				break
			}
		}

		if mv.AudioUrl != "" {
			log.Printf("mv name %s add...\n", mv.Name)
			mvs = append(mvs, mv)
		}
	}

	return mvs
}

func downloadMvToDisk(pathdir string, proxy func(_ *http.Request) (*url.URL, error), mvs ...entity.Mv) {

	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{Transport: transport}
	for _, mv := range mvs {
		response, err := client.Get(mv.AudioUrl)
		if err != nil {
			resp, err1 := client.Get(mv.AudioUrl)
			if err1 != nil {
				continue
			}
			response = resp
		}
		filename := strings.Replace("www.yinyuetai.com_"+mv.Name+"_"+mv.Artist, "/", "", -1)
		filename = strings.Replace(filename, "?", "", -1)
		filename = strings.Replace(filename, "#", "", -1)
		filename = strings.Replace(filename, "|", "", -1)
		filename = strings.Replace(filename, "*", "", -1)
		filename = strings.Replace(filename, ">", "", -1)
		filename = strings.Replace(filename, "<", "", -1)
		filename = strings.Replace(filename, " ", "", -1)
		fileSuffix := mv.AudioUrl[strings.LastIndex(mv.AudioUrl, "."):]

		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			cont, err1 := ioutil.ReadAll(response.Body)
			if err1 != nil {
				log.Printf("Error, %s download failed. skipped.\n", filename)
				continue
			}
			content = cont
		}
		ioutil.WriteFile(pathdir+filename+fileSuffix, content, 0666)
		fmt.Printf("downloading mv %s ...\n", filename)
	}
}


func CrawlOld() {
	//生成组合标签
	seedURLs := getSeedUrls()

	//提取mv id列表
	mvJsonUrls := getMvInfoUrls(50, seedURLs...)
	log.Printf("find mv info request url total : %d \n", len(mvJsonUrls))

	//获取mv的name,artist,audioUrl
	mvs := getMVs(mvJsonUrls...)

	log.Printf("mvs size: %d\n", len(mvs))

	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://127.0.0.1:8087")
	}
	downloadMvToDisk("/Users/sypeng/datas/gomv/", proxy, mvs...)
}