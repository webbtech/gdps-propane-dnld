package xlsx

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/pulpfree/gdps-propane-dwnld/model"
	log "github.com/sirupsen/logrus"
)

// XLSX struct
type XLSX struct {
	file *excelize.File
}

// Defaults
const (
	abc             = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	floatFrmt       = "#,#0"
	timeShortForm   = "20060102"
	timeLongForm    = "2006-01-02"
	timeMonthForm   = "200601"
	dateDayFormat   = "Jan _2"
	dateMonthFormat = "Jan 2006"
	autoFuelID      = "475"
	tankFuelID      = "476"
)

// NewFile function
func NewFile() (x *XLSX, err error) {

	x = new(XLSX)
	x.file = excelize.NewFile()
	if err != nil {
		log.Errorf("xlsx err %s: ", err)
	}
	return x, err
}

// PropaneSales method
func (x *XLSX) PropaneSales(sales *model.PropaneSales) (err error) {

	var cell string
	var style int

	xlsx := x.file
	sheetNm := "Sheet1"

	xlsx.SetSheetName(sheetNm, "Propane Report")

	// Merge cells to accommodate title
	endCell := toChar(4) + "1"
	xlsx.MergeCell(sheetNm, "A1", endCell)

	style, _ = xlsx.NewStyle(`{"font":{"bold":true,"size":12}}`)

	title := fmt.Sprintf("Propane Sales & Deliveries - %s", sales.Report.Date.Format(dateMonthFormat))
	xlsx.SetCellValue(sheetNm, "A1", title)
	xlsx.SetCellStyle(sheetNm, "A1", "A1", style)

	// Create second row with headings
	style, _ = xlsx.NewStyle(`{"font":{"bold":true}}`)
	xlsx.SetColWidth(sheetNm, "B", "D", 10)
	xlsx.SetCellValue(sheetNm, "A2", "Date")
	xlsx.SetCellValue(sheetNm, "B2", "Auto Fuel")
	xlsx.SetCellValue(sheetNm, "C2", "Tank Fuel")
	xlsx.SetCellValue(sheetNm, "D2", "Delivery")
	xlsx.SetCellStyle(sheetNm, "A2", "D2", style)

	// Fill in data
	col := 1
	row := 3
	style, _ = xlsx.NewStyle(`{"number_format": 3}`)

	for _, r := range sales.Report.Sales {

		t, _ := time.Parse(timeLongForm, r.Date)
		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetNm, cell, t.Format(dateDayFormat))

		col++
		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetNm, cell, r.Sales[autoFuelID])
		xlsx.SetCellStyle(sheetNm, cell, cell, style)

		col++
		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetNm, cell, r.Sales[tankFuelID])
		xlsx.SetCellStyle(sheetNm, cell, cell, style)

		if sales.Report.Deliveries[r.Date] > 0 {
			col++
			cell = toChar(col) + strconv.Itoa(row)
			xlsx.SetCellValue(sheetNm, cell, sales.Report.Deliveries[r.Date])
			xlsx.SetCellStyle(sheetNm, cell, cell, style)
		}

		col = 1
		row++
	}

	// Create summary row
	style, _ = xlsx.NewStyle(`{"number_format": 3, "font":{"bold":true}}`)

	col++
	cell = toChar(col) + strconv.Itoa(row)
	cellStart := cell
	rangeStr := fmt.Sprintf("SUM(B3:B%d)", row-1)
	xlsx.SetCellFormula(sheetNm, cell, rangeStr)

	col++
	cell = toChar(col) + strconv.Itoa(row)
	rangeStr = fmt.Sprintf("SUM(C3:C%d)", row-1)
	xlsx.SetCellFormula(sheetNm, cell, rangeStr)

	col++
	cell = toChar(col) + strconv.Itoa(row)
	rangeStr = fmt.Sprintf("SUM(D3:D%d)", row-1)
	xlsx.SetCellFormula(sheetNm, cell, rangeStr)

	xlsx.SetCellStyle(sheetNm, cellStart, cell, style)

	return err
}

// OutputFile method
func (x *XLSX) OutputFile() (buf bytes.Buffer, err error) {
	err = x.file.Write(&buf)
	if err != nil {
		log.Errorf("xlsx err: %s", err)
	}
	return buf, err
}

// OutputToDisk method
func (x *XLSX) OutputToDisk(path string) (fp string, err error) {
	err = x.file.SaveAs(path)
	return path, err
}

// ======================== Helper Methods ================================= //

// see: https://stackoverflow.com/questions/36803999/golang-alphabetic-representation-of-a-number
// for a way to map int to letters
func toChar(i int) string {
	return abc[i-1 : i]
}

// Found these function at: https://stackoverflow.com/questions/18390266/how-can-we-truncate-float64-type-to-a-particular-precision-in-golang
// Looks like a good way to deal with precision
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
