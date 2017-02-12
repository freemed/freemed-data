package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/freemed/freemed-data/common"
	http "github.com/jbuchbinder/gosimplehttp"
	"io/ioutil"
	"strings"
)

const (
	NdcZipUrl      = "http://www.accessdata.fda.gov/cder/ndctext.zip"
	NdcProductFile = "product.txt"
)

var (
	SqlOutput = flag.Bool("sql", false, "Output SQL")
	LocalFile = flag.String("local-file", "", "Read from local ZIP source")
)

func main() {
	var file []byte
	var code int
	var err error

	flag.Parse()

	if *LocalFile == "" {
		fmt.Printf("HTTP GET : %s\n", NdcZipUrl)
		code, file, _, err = http.SimpleGet(NdcZipUrl)
		if err != nil {
			panic(err)
		}
		if code > 299 {
			fmt.Printf("HTTP request got result code %d\n", code)
			return
		}
	} else {
		file, err = ioutil.ReadFile(*LocalFile)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("Extract product file %s from archive\n", NdcProductFile)
	contents, err := common.FileFromZipArchive(file, NdcProductFile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decompressed %d bytes from ZIP archive\n", len(contents))
	r := csv.NewReader(strings.NewReader(string(contents)))
	r.Comma = '\t'
	r.LazyQuotes = true
	rec, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	if *SqlOutput {
		fmt.Printf(SqlPreamble)
		common.InsertsFromArrays("ndc", []string{
			"ProductID",
			"ProductNDC",
			"ProductTypeName",
			"ProprietaryName",
			"ProprietaryNameSuffix",
			"NonProprietaryName",
			"DosageFormName",
			"RouteName",
			"StartMarketingDate",
			"EndMarketingDate",
			"MarketingCategoryName",
			"ApplicationNumber",
			"LabelerName",
			"SubstanceName",
			"StrengthNumber",
			"StrengthUnit",
			"PharmClasses",
			"DEASchedule",
		}, rec)
		return
	}

	//fmt.Printf("%#v", rec[1:])
}
