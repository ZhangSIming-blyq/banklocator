package bankmap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Address 地址结构体
type Address struct {
	Latitude  float64 `json:"lat"` // 纬度
	Longitude float64 `json:"lng"` // 经度
}

// RegeocodeResult 逆地理编码结果结构体
type RegeocodeResult struct {
	AddressComponent struct {
		Province string `json:"province"` // 省份
		City     string `json:"city"`     // 城市
		District string `json:"district"` // 区县
		Street   string `json:"street"`   // 街道名称
	} `json:"addressComponent"`
	Location Address `json:"location"` // 地理位置
}

// GetLocation 调用高德地图逆地理编码 API 获取指定地点的经纬度坐标
func GetLocation(address string) (float64, float64) {
	resp, err := http.Get(fmt.Sprintf("https://restapi.amap.com/v3/geocode/geo?key=%s&address=%s", AmapAPIKey, address))
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	var result struct {
		Status   string `json:"status"`
		Geocodes []struct {
			Location string `json:"location"` // 经纬度坐标
		} `json:"geocodes"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Error(err)
	}
	// return because get a unexpected result
	if result.Status != "1" {
		log.Error("Get empty result, please try again")
		return 0, 0
	}
	if len(result.Geocodes) > 0 {
		// 解析经纬度坐标
		loc := result.Geocodes[0].Location
		var lat, lng float64
		fmt.Sscanf(loc, "%f,%f", &lng, &lat)
		return lat, lng
	} else {
		log.Error("failed to get location for address: %s", address)
	}
	return 0, 0
}

// GenFormatLocation 格式化经纬度输出，经度,维度
func GenFormatLocation(lo string) string {
	// 查询北京市海淀区上地十街10号的经纬度坐标
	lat, lng := GetLocation(lo)
	if lat == 0 && lng == 0 {
		return "get empty result"
	} else {
		return fmt.Sprintf("%.6f,%.6f", lng, lat)
	}
}
