package bankmap

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var (
	Lock              = sync.Mutex{}
	ExportId          = 1
	ExportBankInfoMap = map[int][]BankInfo{}
	KeyWord           string
	Distance          int
	AmapAPIKey        string
)

func init() {
	// KeyWord, default to 银行, try to enter categorical words, such as 学校，餐厅，游乐园
	KeyWord = os.Getenv("KEYWORD")
	if len(KeyWord) == 0 {
		KeyWord = "银行"
	}

	// Distance, the search radius for the given coordinates, default 1000, unit is "meter"
	DistanceS := os.Getenv("DISTANCE")
	if len(DistanceS) == 0 {
		Distance = 1000
	} else {
		Dist, err := strconv.Atoi(DistanceS)
		if err != nil {
			panic("Error in parsing Distance")
		}
		Distance = Dist
	}
	// AmapKey, apikey of Amap developer platform, default ""
	AmapAPIKey = os.Getenv("AMAPKEY")
	if len(AmapAPIKey) == 0 {
		AmapAPIKey = ""
		panic("Amapkey must set!")
	}
}

// GetNearByBank get information about nearby banks and return rendered html
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
	bankList := GetBankList(poiList[0].Location, Distance)
	for i, bank := range bankList {
		bankList[i].Tel = GetTelDetail(bank.ID)
		CalculateScore(&bankList[i])
	}

	sort.Slice(bankList, func(i, j int) bool {
		return bankList[i].Score > bankList[j].Score
	})

	return bankList
}

// GetPOIList get a list of the nearest POIs for a given location
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

// GetBankList get a list of the nearest banks in a given location
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

	url := fmt.Sprintf("%s?keywords=%s&location=%.6f,%.6f&radius=%d&output=json&key=%s&page=%d", AmapAPIURL, KeyWord, latf, lngf, radius, AmapAPIKey, 1)
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
		url = fmt.Sprintf("%s?keywords=%s&location=%.6f,%.6f&radius=%d&output=json&key=%s&page=%d", AmapAPIURL, KeyWord, latf, lngf, radius, AmapAPIKey, ct)
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

// distance calculate the distance between two geographical locations in kilometers
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

// rad converting angles to radians
func rad(d float64) float64 {
	return d * math.Pi / 180.0
}
