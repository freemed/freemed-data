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
	Debug     = flag.Bool("debug", false, "Show debugging information")
	SqlOutput = flag.Bool("sql", false, "Output SQL")
	LocalFile = flag.String("local-file", "", "Read from local ZIP source")
)

func main() {
	var file []byte
	var code int
	var err error

	flag.Parse()

	if *LocalFile == "" {
		if *Debug {
			fmt.Printf("## HTTP GET : %s\n", NdcZipUrl)
		}
		code, file, _, err = http.SimpleGet(NdcZipUrl)
		if err != nil {
			panic(err)
		}
		if code > 299 {
			fmt.Printf("## HTTP request got result code %d\n", code)
			return
		}
	} else {
		file, err = ioutil.ReadFile(*LocalFile)
		if err != nil {
			panic(err)
		}
	}

	if *Debug {
		fmt.Printf("## Extract product file %s from archive\n", NdcProductFile)
	}
	contents, err := common.FileFromZipArchive(file, NdcProductFile, *Debug)
	if err != nil {
		panic(err)
	}
	if *Debug {
		fmt.Printf("## Decompressed %d bytes from ZIP archive\n", len(contents))
	}
	r := csv.NewReader(strings.NewReader(string(contents)))
	r.Comma = '\t'
	r.LazyQuotes = true
	rec, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	// Adjust all date/time values
	for idx, _ := range rec {
		if rec[idx][8] == "" {
			rec[idx][8] = "NULL"
		} else {
			x := rec[idx][8]
			rec[idx][8] = x[0:4] + "-" + x[4:6] + "-" + x[6:8]
		}
		if rec[idx][9] == "" {
			rec[idx][9] = "NULL"
		} else {
			x := rec[idx][9]
			rec[idx][9] = x[0:4] + "-" + x[4:6] + "-" + x[6:8]
		}
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
			"DosageFormName",     // [6]
			"RouteName",          // [7]
			"StartMarketingDate", // [8]
			"EndMarketingDate",   // [9]
			"MarketingCategoryName",
			"ApplicationNumber",
			"LabelerName",
			"SubstanceName",
			"StrengthNumber",
			"StrengthUnit", // [15]
			"PharmClasses", // [16]
			"DEASchedule",
		}, rec[1:])
		common.InsertsFromArrays("ndcDosageForm", []string{
			"DosageFormName",
		}, common.OneToMultiArray(Derivatives(rec[1:], 6, ";")))
		common.InsertsFromArrays("ndcRoute", []string{
			"RouteName",
		}, common.OneToMultiArray(Derivatives(rec[1:], 7, ";")))
		common.InsertsFromArrays("ndcStrengthUnit", []string{
			"StrengthUnit",
		}, common.OneToMultiArray(Derivatives(rec[1:], 15, ";")))
		common.InsertsFromArrays("ndcPharmClass", []string{
			"PharmClassName",
		}, common.OneToMultiArray(Derivatives(rec[1:], 16, ";")))
		return
	}

	fmt.Printf("DosageForm : %#v\n", Derivatives(rec[1:], 6, ";"))
	fmt.Printf("Route : %#v\n", Derivatives(rec[1:], 7, ";"))
	fmt.Printf("StrengthUnit : %#v\n", Derivatives(rec[1:], 15, ";"))
	fmt.Printf("PharmClasses : %#v\n", Derivatives(rec[1:], 16, ","))
	//fmt.Printf("%#v", rec[1:])
}
