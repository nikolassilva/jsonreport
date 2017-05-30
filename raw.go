package jsonreport

import (
	"encoding/json"
)

// Raw TODO
func Raw(data []byte) (*Report, error) {
	var rows rawRows
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, err
	}

	headers := rows.headers()
	records := rows.records(headers)
	report := &Report{
		Headers: headers,
		Records: records,
	}

	return report, nil
}

type rawRow map[string]json.RawMessage

func (row rawRow) record(headers []string) []string {
	record := make([]string, len(headers))

	for i, header := range headers {
		v, ok := row[header]

		if ok {
			record[i] = formatJSONValue(v)
		} else {
			record[i] = ""
		}
	}

	return record
}

func formatJSONValue(v []byte) string {
	// checks if first character is a double quote
	if v[0] == 0x22 {
		var str string
		json.Unmarshal(v, &str)
		return str
	}

	str := string(v)
	if str == "null" {
		return ""
	}

	return str
}

type rawRows []rawRow

func (rows rawRows) headers() []string {
	firstRow := rows[0]
	headers := make([]string, len(firstRow))

	i := 0
	for key := range firstRow {
		headers[i] = key
		i++
	}

	return headers
}

func (rows rawRows) records(headers []string) [][]string {
	records := make([][]string, len(rows))

	i := 0
	for _, row := range rows {
		records[i] = row.record(headers)
		i++
	}

	return records
}
