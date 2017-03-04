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
	TsvMode   = flag.Bool("tsv", false, "TSV build/refresh mode")
	LocalFile = flag.String("local-file", "", "Read from local ZIP source")
)

func main() {
	var file []byte
	var code int
	var err error

	flag.Parse()

	if *TsvMode && *SqlOutput {
		flag.PrintDefaults()
		return
	}

	// ----- INGEST DATA -----

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

	// ----- EXPORT DATA -----

	if *TsvMode {
		fmt.Println("* Using TSV Mode")

		// Determine whether or not we're updating or not ...
		updateMode := false
		if common.FileExists("data/route.tsv") {
			updateMode = true
			fmt.Println(" * Detected update mode")
		}

		if updateMode {
			// Update mode :

			fmt.Println(" - Rewriting new ndc.tsv")
			err := common.TsvFromArrays(
				"data/ndc.tsv",
				common.PrependUniqueIds(
					rec,
				),
			)
			if err != nil {
				panic(err)
			}

			// route.tsv
			{
				fmt.Println(" - Ingesting route.tsv")
				raw, err := common.ReadTsv("data/route.tsv")
				if err != nil {
					panic(err)
				}
				newData := false
				corpus := common.Derivatives(raw, 1, ";")
				updates := common.Derivatives(rec[1:], 7, ";")
				keys := common.CoerceSliceStringToInt(common.Derivatives(raw, 0, ";"))
				maxVal := common.MaxIntSlice(keys)
				for _, u := range updates {
					uc := strings.TrimSpace(strings.ToUpper(u))
					if uc == "" {
						continue
					}
					if !common.HasElement(corpus, uc) {
						newData = true
						maxVal++
						fmt.Printf(" ! Found new element '%s' (%d)\n", uc, maxVal)
						raw = append(raw, []string{fmt.Sprintf("%d", maxVal), uc})
						corpus = append(corpus, uc)
					}
				}

				// Push out to TSV
				if newData {
					fmt.Println(" - Writing updated route.tsv")
					err = common.TsvFromArrays(
						"data/route.tsv",
						raw,
					)
					if err != nil {
						panic(err)
					}
				} else {
					fmt.Println(" - No new data to write to routes.tsv")
				}
			}
		} else {
			// Create mode : blast data out to files

			fmt.Println(" - Writing new ndc.tsv")
			err := common.TsvFromArrays(
				"data/ndc.tsv",
				common.PrependUniqueIds(
					rec,
				),
			)
			if err != nil {
				panic(err)
			}
			fmt.Println(" - Writing new route.tsv")
			err = common.TsvFromArrays(
				"data/route.tsv",
				common.PrependUniqueIds(
					common.OneToMultiArray(common.Derivatives(rec[1:], 7, ";"), true),
				),
			)
			if err != nil {
				panic(err)
			}

			fmt.Println(" - Writing new producttype.tsv")
			err = common.TsvFromArrays(
				"data/producttype.tsv",
				common.PrependUniqueIds(
					common.OneToMultiArray(common.Derivatives(rec[1:], 2, ";"), true),
				),
			)
			if err != nil {
				panic(err)
			}

			fmt.Println(" - Writing new dosageform.tsv")
			err = common.TsvFromArrays(
				"data/dosageform.tsv",
				common.PrependUniqueIds(
					common.OneToMultiArray(common.Derivatives(rec[1:], 6, ";"), true),
				),
			)
			if err != nil {
				panic(err)
			}

			fmt.Println(" - Writing new strengthunit.tsv")
			err = common.TsvFromArrays(
				"data/strengthunit.tsv",
				common.PrependUniqueIds(
					common.OneToMultiArray(common.Derivatives(rec[1:], 15, ";"), true),
				),
			)
			if err != nil {
				panic(err)
			}

			fmt.Println(" - Writing new pharmclass.tsv")
			err = common.TsvFromArrays(
				"data/pharmclass.tsv",
				common.PrependUniqueIds(
					common.OneToMultiArray(common.Derivatives(rec[1:], 16, ","), true),
				),
			)
			if err != nil {
				panic(err)
			}

			fmt.Println(" - Writing new drugname.tsv")
			err = common.TsvFromArrays(
				"data/drugname.tsv",
				common.PrependUniqueIds(
					common.OneToMultiArray(common.Derivatives(rec[1:], 3, ";"), true),
				),
			)
			if err != nil {
				panic(err)
			}
		}

		return
	}

	if *SqlOutput {
		fmt.Printf(SqlPreamble)
		common.InsertsFromArrays("ndc", []string{
			"ProductID",
			"ProductNDC",
			"ProductTypeName", // [2]
			"ProprietaryName", // [3]
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
		}, common.OneToMultiArray(common.Derivatives(rec[1:], 6, ";"), true))
		common.InsertsFromArrays("ndcRoute", []string{
			"RouteName",
		}, common.OneToMultiArray(common.Derivatives(rec[1:], 7, ";"), true))
		common.InsertsFromArrays("ndcStrengthUnit", []string{
			"StrengthUnit",
		}, common.OneToMultiArray(common.Derivatives(rec[1:], 15, ";"), true))
		common.InsertsFromArrays("ndcPharmClass", []string{
			"PharmClassName",
		}, common.OneToMultiArray(common.Derivatives(rec[1:], 16, ","), true))
		return
	}

	fmt.Printf("DosageForm : %#v\n", common.Derivatives(rec[1:], 6, ";"))
	fmt.Printf("Route : %#v\n", common.Derivatives(rec[1:], 7, ";"))
	fmt.Printf("StrengthUnit : %#v\n", common.Derivatives(rec[1:], 15, ";"))
	fmt.Printf("PharmClasses : %#v\n", common.Derivatives(rec[1:], 16, ","))
	//fmt.Printf("%#v", rec[1:])
}
