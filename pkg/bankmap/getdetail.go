package bankmap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

// GetTelDetail get tel detail from poi info
func GetTelDetail(id string) string {
	output := "json"

	endpoint := "https://restapi.amap.com/v3/place/detail"
	params := url.Values{}
	params.Set("key", AmapAPIKey)
	params.Set("id", id)
	params.Set("output", output)

	url := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	var rspnse Response
	err = json.Unmarshal(body, &rspnse)
	if err != nil {
		log.Error(err)
	}

	// get merchant tel
	tel := ""
	if len(rspnse.Poi) != 0 {
		for i := range rspnse.Poi {
			if reflect.TypeOf(rspnse.Poi[i].Tel).Kind() == reflect.String {
				tel += rspnse.Poi[i].Tel.(string)
			}
		}
	}
	return tel
}

type Response struct {
	Status string `json:"status"`
	Info   string `json:"info"`
	Poi    []Poi  `json:"pois"`
}

type Poi struct {
	Name string `json:"name"`
	Tel  any    `json:"tel"`
}
