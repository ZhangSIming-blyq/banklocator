package bankmap

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var CountryBankList = []string{"中国建设银行", "中国工商银行", "中国农业银行", "中国邮政储蓄银行", "交通银行", "中国银行"}
var SubBankList = []string{"中国农业发展银行", "中国进出口银行", "中国光大银行", "中国民生银行", "招商银行"}

func CalculateScore(bi *BankInfo) {
	// distance aspect
	dis, err := stof(bi.Distance)
	if err != nil {
		log.Error(dis)
	}
	dtk := float64(Distance / 1000)
	dP := (dtk - dis) / dtk * 0.55
	// priority aspect
	pP := 0.5 * 0.35
	for _, cB := range CountryBankList {
		if strings.Contains(bi.Name, cB) {
			pP = 1 * 0.35
			break
		}
	}

	for _, sB := range SubBankList {
		if strings.Contains(bi.Name, sB) {
			pP = 0.75 * 0.35
			break
		}
	}

	// telephone aspect
	tP := 0.0
	if bi.Tel != "" {
		tP = 1 * 0.1
	}
	totalScore := dP + pP + tP
	totalScoreFmt := fmt.Sprintf("%.2f", totalScore*100)
	bi.Score = totalScoreFmt
}

func stof(strValue string) (float64, error) {
	value, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return 0, err
	}
	return math.Round(value*100) / 100, nil
}
