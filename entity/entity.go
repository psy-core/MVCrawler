package entity

type Mv struct {
	Name     string
	Artist   string
	AudioUrl string
}

type MvInfo struct {
	VideoInfo VideoInfo `json:"videoInfo"`
}

type VideoInfo struct {
	CoreVideoInfo CoreVideoInfo `json:"coreVideoInfo"`
}

type CoreVideoInfo struct {
	ArtistName     string `json:"artistName"`
	VideoName      string `json:"videoName"`
	VideoUrlModels []VideoUrlModels `json:"videoUrlModels"`
}

type VideoUrlModels struct {
	QualityLevel string `json:"qualityLevel"`
	VideoUrl     string `json:"videoUrl"`
}