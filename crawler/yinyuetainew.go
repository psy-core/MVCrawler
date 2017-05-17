package crawler

import (
	"net/http"
	"github.com/cihub/seelog"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"net/url"
	"github.com/psy-core/MVCrawler/entity"
)

func CrawNew() {

	seedURLs := getAPISeedURLs()
	infoURLs := getMVInfoURLsByAPI(seedURLs...)
	seelog.Infof("find mv info request url total : %d \n", len(infoURLs))

	//获取mv的name,artist,audioUrl
	mvs := getMVs(infoURLs...)

	seelog.Infof("mvs size: %d\n", len(mvs))

	//mv去重
	duplicateMap := loadDuplicateMapByDir("E:\\gomv05-09\\")
	mvs = duplicate(duplicateMap, mvs)
	seelog.Infof("mvs size: %d after duplicate.\n", len(mvs))

	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://192.168.86.121:1800")
	}
	downloadMvToDisk("E:\\gomv05-16\\", proxy, mvs...)

}

func getAPISeedURLs() []string {

	return []string{
		"http://mvapi.yinyuetai.com/mvchannel/so?sid=5%3B12&tid=33%3B73&a=&p=&c=sh&s=pubdate&pageSize=100&page=1",
		"http://mvapi.yinyuetai.com/mvchannel/so?sid=5%3B12&tid=30%3B70&a=&p=&c=sh&s=pubdate&pageSize=100&page=1",
		"http://mvapi.yinyuetai.com/mvchannel/so?sid=5%3B12&tid=31%3B71&a=&p=&c=sh&s=pubdate&pageSize=100&page=1",
		"http://mvapi.yinyuetai.com/mvchannel/so?sid=3%3B10&tid=19%3B59&a=&p=&c=sh&s=pubdate&pageSize=100&page=1",
		"http://mvapi.yinyuetai.com/mvchannel/so?sid=3%3B10&tid=17%3B57&a=&p=&c=sh&s=pubdate&pageSize=100&page=1",
	}
}

func getMVInfoURLsByAPI(seedURLs ...string) []string {

	infoURLs := []string{}
	for _, u := range seedURLs {
		resp, err := http.Get(u)
		if err != nil {
			seelog.Error("url ", u, "download failed.", err)
			continue
		}
		jsonbyte, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			seelog.Error("error! http response body read failed. u: ", u)
			continue
		}
		var f interface{}
		err = json.Unmarshal(jsonbyte, &f)
		if err != nil {
			seelog.Error("error! http response json unmarshal failed. u: ", u)
			continue
		}
		m := f.(map[string]interface{})
		result := m["result"]
		mvs := result.([]interface{})
		if len(mvs) == 0 {
			seelog.Warn("json mv list is empty. u:", u)
			continue
		}
		for _, mv := range mvs {
			mvobj := mv.(map[string]interface{})
			videoId := mvobj["videoId"].(float64)
			if videoId > 0 {
				infoURLs = append(infoURLs,
					"http://www.yinyuetai.com/insite/get-video-info?json=true&videoId="+
						strconv.FormatFloat(videoId, 'f', 0, 64))
			}
		}
	}

	return infoURLs
}

func loadDuplicateMapByDir(fileDir string) map[string]bool {

	result := make(map[string]bool)
	fileinfos, err := ioutil.ReadDir(fileDir)
	if err != nil {
		return result
	}
	for _, info := range fileinfos {
		result[info.Name()] = true
	}
	return result
}

func duplicate(duplicate map[string]bool, mvs []entity.Mv) []entity.Mv {

	result := make([]entity.Mv, 0)
	for _, mv := range mvs {
		filename := generateFileName(mv.Name, mv.Artist, mv.AudioUrl)
		if !duplicate[filename] {
			result = append(result, mv)
		} else {
			seelog.Warnf("file %s is duplicate. ignored.", filename)
		}
	}
	return result
}
