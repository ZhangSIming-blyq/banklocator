package bankmap

import "testing"

func TestCalculateScore(t *testing.T) {
	bi := BankInfo{
		Name:     "aslkj",
		Distance: "0.1",
	}
	CalculateScore(&bi)
	t.Log(bi.Score)
}
