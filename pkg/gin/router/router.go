package gin

import (
	"banklocator/pkg/bankmap"
	"banklocator/pkg/gin/middleware"
	"banklocator/pkg/logger"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var log = logger.DefaultLog

func InitRouter() {
	router := gin.Default()

	// Use the AuthMiddleware for all routes
	app := router.Group("/app")
	app.Use(middleware.AuthMiddleware())

	router.LoadHTMLGlob("templates/*")
	router.Static("/public", "./static")

	router.GET("/", func(c *gin.Context) {

		if c.Writer.Status() >= 300 && c.Writer.Status() < 400 {
			c.HTML(http.StatusTemporaryRedirect, "index.html", gin.H{})
		} else {
			c.HTML(http.StatusOK, "index.html", gin.H{})
		}
	})

	// for test
	router.GET("/test", func(c *gin.Context) {
		a := bankmap.BankInfo{
			Name:     "test",
			Distance: "0.1",
			Tel:      "test",
			Score:    "93.1",
		}

		b := bankmap.BankInfo{
			Name:     "test",
			Distance: "0.1",
			Tel:      "test",
			Score:    "93.1",
		}

		d := bankmap.BankInfo{
			Name:     "test",
			Distance: "0.1",
			Tel:      "test",
			Score:    "93.1",
		}

		var rst []bankmap.BankInfo
		rst = append(rst, a)
		rst = append(rst, b)
		rst = append(rst, d)

		bankmap.Lock.Lock()
		bankmap.ExportId += 1
		bankmap.ExportBankInfoMap[bankmap.ExportId] = rst
		c.HTML(http.StatusOK, "result.html", gin.H{
			"name":     rst,
			"exportid": bankmap.ExportId,
		})
		bankmap.Lock.Unlock()
	})

	app.POST("/submit", func(c *gin.Context) {

		var rstList [][]bankmap.BankInfo
		name := c.PostForm("sname")
		if name == "" {
			c.HTML(http.StatusOK, "result.html", gin.H{})
			return
		}
		nameList := []string{name}
		if strings.Contains(name, ",") {
			nameList = strings.Split(name, ",")
		} else if strings.Contains(name, "，") {
			nameList = strings.Split(name, "，")
		}
		for _, name := range nameList {
			rst := bankmap.GetNearByBank(name)
			rstList = append(rstList, rst)
		}

		if len(rstList) == 1 {
			if rstList[0] == nil {
				c.HTML(http.StatusOK, "result.html", gin.H{})
				return
			}
			bankmap.Lock.Lock()
			bankmap.ExportId += 1
			bankmap.ExportBankInfoMap[bankmap.ExportId] = rstList[0]
			c.HTML(http.StatusOK, "result.html", gin.H{
				"name":     rstList[0],
				"exportid": bankmap.ExportId,
			})
			bankmap.Lock.Unlock()
		} else {
			var interSlice []bankmap.BankInfo
			for i := range rstList {
				if i == len(rstList)-1 {
					break
				}
				interSlice = rstList[i]
				interSlice = Intersect(interSlice, rstList[i+1])
			}
			sort.Slice(interSlice, func(i, j int) bool {
				return interSlice[i].Score > interSlice[j].Score
			})
			bankmap.Lock.Lock()
			bankmap.ExportId += 1
			bankmap.ExportBankInfoMap[bankmap.ExportId] = interSlice
			c.HTML(http.StatusOK, "result.html", gin.H{
				"name":     interSlice,
				"exportid": bankmap.ExportId,
			})
			bankmap.Lock.Unlock()
		}

	})
	app.GET("/submit", func(c *gin.Context) {
		c.HTML(http.StatusOK, "result.html", gin.H{})
	})

	// for export
	app.GET("/export", middleware.ExportHandler)

	router.Run(":30080")
}

func Intersect(slice1, slice2 []bankmap.BankInfo) []bankmap.BankInfo {
	nameDisMap := make(map[string]string)
	set := make(map[string]bool)
	var intersect []bankmap.BankInfo

	for _, value := range slice1 {
		set[value.Name] = true
		nameDisMap[value.Name] = value.Distance
	}

	for _, value := range slice2 {
		if set[value.Name] {
			// Average distance
			ndis, err := stof(value.Distance)
			if err != nil {
				log.Error(err)
			}
			odis, err := stof(nameDisMap[value.Name])
			if err != nil {
				log.Error(err)
			}
			ad := (ndis + odis) / 2
			value.Distance = fmt.Sprintf("%.2f", ad)
			bankmap.CalculateScore(&value)
			intersect = append(intersect, value)
		}
	}

	return intersect
}

func stof(strValue string) (float64, error) {
	value, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return 0, err
	}
	return math.Round(value*100) / 100, nil
}
