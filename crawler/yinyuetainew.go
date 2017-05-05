package crawler

import (
	"net/http"
	"github.com/cihub/seelog"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"net/url"
)

func CrawNew() {

	seedURLs := getAPISeedURLs()
	infoURLs := getMVInfoURLsByAPI(seedURLs...)
	seelog.Info("find mv info request url total : %d \n", len(infoURLs))

	//获取mv的name,artist,audioUrl
	mvs := getMVs(infoURLs...)

	seelog.Info("mvs size: %d\n", len(mvs))

	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://127.0.0.1:8087")
	}
	downloadMvToDisk("/Users/sypeng/datas/gomv/", proxy, mvs...)

}

func getAPISeedURLs() []string {

	return []string{
		"http://mvapi.yinyuetai.com/mvchannel/so?sid=5%3B12&tid=33%3B73&a=&p=&c=sh&s=pubdate&pageSize=10&page=1",
		"http://mvapi.yinyuetai.com/mvchannel/so?sid=5%3B12&tid=30%3B70&a=&p=&c=sh&s=pubdate&pageSize=10&page=1",
		"http://mvapi.yinyuetai.com/mvchannel/so?sid=5%3B12&tid=31%3B71&a=&p=&c=sh&s=pubdate&pageSize=10&page=1",
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
