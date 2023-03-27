package bankmap

import "banklocator/pkg/logger"

// AmapAPIURL Amap's API URL
const AmapAPIURL = "https://restapi.amap.com/v3/place/around"

var log = logger.DefaultLog

// PlaceResult the result of location query
type PlaceResult struct {
	Status  string        `json:"status"`
	Info    string        `json:"info"`
	Count   string        `json:"count"`
	PoiList []AmapPOIInfo `json:"pois"`
}

// AmapPOIInfo Amap's POI info
type AmapPOIInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Distance string `json:"distance"`
}

// BankInfo the info of bank
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
