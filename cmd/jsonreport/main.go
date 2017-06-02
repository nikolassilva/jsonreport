package main

import (
	"log"
	"os"

	docopt "github.com/nikolassilva/docopt.go"
	"github.com/nikolassilva/jsonreport"
)

const version = "JSONReport 1.0.0"
const usage = `
  Usage:
    jsonreport [--format=] [<json>]

  Options:
    -f --format=FORMAT   Specify the report format [default: csv].
    -h --help            Show this.
    -V --version         Show version.

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

	jsonStream := args["<json>"].(string)
	report, err := jsonreport.Raw([]byte(jsonStream))
	if err != nil {
		log.Fatalf("error parsing json: %s", err)
	}

	switch args["--format"] {
	case "csv":
		if err := report.WriteCSV(os.Stdout); err != nil {
			log.Fatalf("error generating report: %s", err)
		}
	case "xlsx", "excel":
		if err := report.WriteXLSX(os.Stdout); err != nil {
			log.Fatalf("error generating report: %s", err)
		}
	default:
		log.Fatalf("invalid report format. use one of these: csv, xlsx")
	}
}
