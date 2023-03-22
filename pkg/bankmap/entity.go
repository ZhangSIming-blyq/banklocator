package bankmap

import "banklocator/pkg/logger"

// AmapAPIURL 高德地图 API URL
const AmapAPIURL = "https://restapi.amap.com/v3/place/around"

const AmapAPIKey = "b87b49eab52a5769a8973a7e7b98afaa"

var log = logger.DefaultLog

// PlaceResult 地点查询结果
type PlaceResult struct {
	Status  string        `json:"status"`
	Info    string        `json:"info"`
	Count   string        `json:"count"`
	PoiList []AmapPOIInfo `json:"pois"`
}

// AmapPOIInfo 高德地图 POI 信息
type AmapPOIInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Distance string `json:"distance"`
}

// BankInfo 银行信息
type BankInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Distance string `json:"distance"`
	Tel      string `json:"tel"`
	Score    string `json:"score"`
}

type SearchResult struct {
	Results []string
}

type ExportBankInfo struct {
	ID           int
	BankInfoList []BankInfo
}
