package jsonreport

import (
	"encoding/csv"
	"io"
)

// Report TODO
type Report struct {
	Headers []string
	Records [][]string
	Options Options
}

// Options TODO
type Options struct {
	Comma   rune
	UseCRLF bool
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
