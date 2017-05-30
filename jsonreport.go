package jsonreport

import (
	"encoding/csv"
	"io"
	"strconv"

	"github.com/Luxurioust/excelize"
)

// Report TODO
type Report struct {
	Headers []string
	Records [][]string
	Options Options
}

// Options TODO
type Options struct {
	Comma     rune
	UseCRLF   bool
	SheetName string
}

// WriteCSV TODO
func (report *Report) WriteCSV(w io.Writer) error {
	writer := csv.NewWriter(w)

	if report.Options.Comma != 0 {
		writer.Comma = report.Options.Comma
	}
	if report.Options.UseCRLF {
		writer.UseCRLF = report.Options.UseCRLF
	}
	if err := writer.Write(report.Headers); err != nil {
		return err
	}

	return writer.WriteAll(report.Records)
}

// WriteXLSX TODO
func (report *Report) WriteXLSX(w io.Writer) error {
	xlsx := excelize.CreateFile()
	sheetName := xlsx.GetSheetName(xlsx.GetActiveSheetIndex())

	setExcelRow(xlsx, sheetName, report.Headers, 0)
	for i, row := range report.Records {
		setExcelRow(xlsx, sheetName, row, i+1)
	}

	if report.Options.SheetName != "" {
		xlsx.SetSheetName(sheetName, report.Options.SheetName)
	}

	return xlsx.Write(w)
}

func setExcelRow(xlsx *excelize.File, sheetName string, columns []string, rowIndex int) {
	for i, column := range columns {
		cell := colStr(i+1) + strconv.Itoa(rowIndex+1)
		xlsx.SetCellStr(sheetName, cell, column)
	}
}

func colStr(i int) string {
	var str string
	for i > 0 {
		rest := (i - 1) % 26
		str = string(rest+65) + str
		i = (i - rest) / 26
	}

	return str
}
