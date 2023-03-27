package bankmap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Address struct {
	Latitude  float64 `json:"lat"` // 纬度
	Longitude float64 `json:"lng"` // 经度
}

// RegeocodeResult inverse geocoding result structure
type RegeocodeResult struct {
	AddressComponent struct {
		Province string `json:"province"`
		City     string `json:"city"`
		District string `json:"district"`
		Street   string `json:"street"`
	} `json:"addressComponent"`
	Location Address `json:"location"`
}

// GetLocation call the inverse geocoding API of  Amap to get the latitude and longitude coordinates of the specified location
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
			Location string `json:"location"`
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
		// resolve latitude and longitude coordinates
		loc := result.Geocodes[0].Location
		var lat, lng float64
		fmt.Sscanf(loc, "%f,%f", &lng, &lat)
		return lat, lng
	} else {
		log.Error("failed to get location for address: %s", address)
	}
	return 0, 0
}

// GenFormatLocation formatted latitude and longitude output, longitude,dimension
func GenFormatLocation(lo string) string {
	lat, lng := GetLocation(lo)
	if lat == 0 && lng == 0 {
		return "get empty result"
	} else {
		return fmt.Sprintf("%.6f,%.6f", lng, lat)
	}
}
