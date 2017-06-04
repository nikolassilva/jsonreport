package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	docopt "github.com/nikolassilva/docopt.go"
	"github.com/nikolassilva/jsonreport"
)

const version = "JSONReport 1.0.0"
const usage = `
  Usage:
    jsonreport [options] [<json>]

  Options:
    -f, --format=FORMAT    Specify the report format [default: csv].
    -O, --output=FILENAME  Specify the file to write the output.
    -h, --help             Show this.
    -V, --version          Show version.

  Examples:
    pass json as argument and output as excel spreadsheet
    $ jsonreport --format=xlsx '[{"a": 1, "b": 2, "c": 3}, {"a": 4, "b": 5, "c": 6}]'
    
    get json from stdin and output as csv document
    $ cat data.json | jsonreport --format=csv
`

func main() {
	args, err := docopt.Parse(usage, nil, true, version, true)
	if err != nil {
		log.Fatalf("error parsing arguments: %s", err)
	}

	report, err := generateReport(args["<json>"])
	if err != nil {
		log.Fatal(err)
	}

	file, err := getOutputFile(args["--output"])
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	switch args["--format"] {
	case "csv":
		if err := report.WriteCSV(file); err != nil {
			log.Fatalf("error generating report: %s", err)
		}
	case "xlsx", "excel":
		if err := report.WriteXLSX(file); err != nil {
			log.Fatalf("error generating report: %s", err)
		}
	default:
		log.Fatalf("invalid report format. use one of these: csv, xlsx")
	}
}

func generateReport(jsonStream interface{}) (*jsonreport.Report, error) {
	switch v := jsonStream.(type) {
	case string:
		return reportFromJSON([]byte(v))

	default:
		fi, err := os.Stdin.Stat()
		if err != nil {
			return nil, err
		} else if fi.Mode()&os.ModeNamedPipe == 0 {
			return nil, fmt.Errorf("expected pipe or argument but got nothing")
		}

		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("error while reading from stdin: %s", err)
		}

		return reportFromJSON(b)
	}
}

func reportFromJSON(jsonStream []byte) (*jsonreport.Report, error) {
	report, err := jsonreport.Raw([]byte(jsonStream))
	if err != nil {
		return nil, fmt.Errorf("error parsing json: %s", err)
	}

	return report, nil
}

func getOutputFile(filename interface{}) (*os.File, error) {
	switch v := filename.(type) {
	case string:
		file, err := os.Create(v)
		if err != nil {
			return nil, fmt.Errorf("error while opening file '%s': %s", v, err)
		}

		return file, nil
	default:
		return os.Stdout, nil
	}
}
