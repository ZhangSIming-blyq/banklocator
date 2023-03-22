package bankmap

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var (
	Lock              = sync.Mutex{}
	ExportId          = 1
	ExportBankInfoMap = map[int][]BankInfo{}
)

// TODO: achieve this
func init() {
	// keyword
	// distance
	// ampkey
}

// GetNearByBank 获取附近的银行信息，并返回渲染html
func GetNearByBank(name string) []BankInfo {
	var totalMsg SearchResult
	// 1. get latitude and longitude for the address you provided, result count is 1, return 0,0 if met some error
	lo := GenFormatLocation(name)
	if lo == "get empty result" {
		totalMsg.Results = append(totalMsg.Results, "")
		// return totalMsg.Results
		return []BankInfo{}
	}

	// 2. get all merchant which is located within 5000m range
	poiList, err := GetPOIList(name, 5000, lo)
	if err != nil {
		log.Error(err)
		return nil
	}

	// 3. get all bank which is located within 5000m range
	bankList := GetBankList(poiList[0].Location, 1000)
	for i, bank := range bankList {
		bankList[i].Tel = GetTelDetail(bank.ID)
		CalculateScore(&bankList[i])
	}

	sort.Slice(bankList, func(i, j int) bool {
		return bankList[i].Score > bankList[j].Score
	})

	return bankList
}

// GetPOIList 获取指定地点最近的 POI 列表
func GetPOIList(placeName string, radius int, lo string) ([]AmapPOIInfo, error) {
	resp, err := http.Get(fmt.Sprintf("%s?location=%s&radius=%d&output=json&key=%s", AmapAPIURL, lo, radius, AmapAPIKey))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result PlaceResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	if result.Status != "1" {
		return nil, fmt.Errorf(result.Info)
	}
	return result.PoiList, nil
}

// GetBankList 获取指定地点最近的银行列表
func GetBankList(lo string, radius int) []BankInfo {
	losp := strings.Split(lo, ",")
	lat := losp[1]
	lng := losp[0]
	latf, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		log.Error(err)
	}
	lngf, err := strconv.ParseFloat(lng, 64)
	if err != nil {
		log.Error(err)
	}

	bankList := []BankInfo{}
	RstCount := 0
	ct := 1

	url := fmt.Sprintf("%s?keywords=银行&location=%.6f,%.6f&radius=%d&output=json&key=%s&page=%d", AmapAPIURL, latf, lngf, radius, AmapAPIKey, 1)
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return []BankInfo{}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	var result PlaceResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Error(err)
	}
	for _, poi := range result.PoiList {
		// check if need to remove self-service cash machines
		if strings.Contains(poi.Name, "自助") || strings.Contains(poi.Name, "ATM") || strings.Contains(poi.Name, "暂停") {
			continue
		}
		tmpBI := BankInfo{
			ID:       poi.ID,
			Name:     poi.Name,
			Location: fmt.Sprintf(poi.Location),
			// poi.Location is corporate's location and latf/lngf is bank's coordinate
			Distance: fmt.Sprintf("%.2f", distance(latf, lngf, poi.Location)),
		}
		bankList = append(bankList, tmpBI)
	}

	num, err := strconv.Atoi(result.Count)
	if err != nil {
		log.Error(err)
	}
	RstCount = num
	var tc int
	if RstCount/20 > 5 {
		tc = 5
	} else {
		tc = RstCount / 20
	}

	for ct = 2; ct < tc; ct++ {
		url = fmt.Sprintf("%s?keywords=银行&location=%.6f,%.6f&radius=%d&output=json&key=%s&page=%d", AmapAPIURL, latf, lngf, radius, AmapAPIKey, ct)
		resp, err := http.Get(url)
		if err != nil {
			log.Error(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
		}
		var result PlaceResult
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Error(err)
		}

		for _, poi := range result.PoiList {
			// check if need to remove self-service cash machines
			if strings.Contains(poi.Name, "自助") || strings.Contains(poi.Name, "ATM") || strings.Contains(poi.Name, "暂停") {
				continue
			}
			tmpBI := BankInfo{
				ID:       poi.ID,
				Name:     poi.Name,
				Location: fmt.Sprintf(poi.Location),
				// poi.Location is corporate's location and latf/lngf is bank's coordinate
				Distance: fmt.Sprintf("%.2f", distance(latf, lngf, poi.Location)),
			}
			bankList = append(bankList, tmpBI)
		}
	}
	log.Info("Total ct count is ", ct)

	return bankList
}

// distance 计算两个地理位置之间的距离，单位为千米
func distance(lat1, lng1 float64, poiLo string) float64 {
	poiLoSp := strings.Split(poiLo, ",")
	lat2 := poiLoSp[1]
	lng2 := poiLoSp[0]
	latf, err := strconv.ParseFloat(lat2, 64)
	if err != nil {
		log.Error(err)
	}
	lngf, err := strconv.ParseFloat(lng2, 64)
	if err != nil {
		log.Error(err)
	}
	radLat1 := rad(lat1)
	radLat2 := rad(latf)
	a := radLat1 - radLat2
	b := rad(lng1) - rad(lngf)
	s := 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(a/2), 2)+math.Cos(radLat1)*math.Cos(radLat2)*math.Pow(math.Sin(b/2), 2)))
	s = s * 6378.137
	return s
}

// rad 将角度转化为弧度
func rad(d float64) float64 {
	return d * math.Pi / 180.0
}
