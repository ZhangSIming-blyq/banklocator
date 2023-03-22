package middleware

import (
	"banklocator/pkg/bankmap"
	"banklocator/pkg/logger"
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"strconv"
)

var log = logger.DefaultLog

func ExportHandler(c *gin.Context) {
	eid := c.Query("id")

	num, err := strconv.Atoi(eid)
	if err != nil {
		log.Error(err)
		return
	}

	// Get the data to export
	bankList := bankmap.ExportBankInfoMap[num]

	// Set the response headers to indicate that this is a CSV file
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment;filename=banklist.csv")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Create a new CSV writer
	w := csv.NewWriter(c.Writer)

	// Write the CSV header row
	w.Write([]string{"Name", "Distance", "Tel", "Score"})

	// Write each row of data
	for _, bank := range bankList {
		row := []string{bank.Name, bank.Distance, bank.Tel, bank.Score}
		w.Write(row)
	}

	// Flush the writer to ensure all data is written to the response
	w.Flush()

}
