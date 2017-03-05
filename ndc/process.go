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

			mergeTsv("route.tsv", rec, 7, ";")
			mergeTsv("producttype.tsv", rec, 2, ";")
			mergeTsv("dosageform.tsv", rec, 6, ";")
			mergeTsv("strengthunit.tsv", rec, 15, ";")
			mergeTsv("pharmclass.tsv", rec, 16, ",")
			mergeTsv("drugname.tsv", rec, 3, ";")
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

			newTsv("route.tsv", rec, 7, ";")
			newTsv("producttype.tsv", rec, 2, ";")
			newTsv("dosageform.tsv", rec, 6, ";")
			newTsv("strengthunit.tsv", rec, 15, ";")
			newTsv("pharmclass.tsv", rec, 16, ",")
			newTsv("drugname.tsv", rec, 3, ";")
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
		common.InsertsFromArrays("ndcRoute", []string{
			"RouteName",
		}, common.OneToMultiArray(common.Derivatives(rec[1:], 7, ";"), true))
		common.InsertsFromArrays("ndcProductType", []string{
			"ProductType",
		}, common.OneToMultiArray(common.Derivatives(rec[1:], 2, ";"), true))
		common.InsertsFromArrays("ndcDosageForm", []string{
			"DosageFormName",
		}, common.OneToMultiArray(common.Derivatives(rec[1:], 6, ";"), true))
		common.InsertsFromArrays("ndcStrengthUnit", []string{
			"StrengthUnit",
		}, common.OneToMultiArray(common.Derivatives(rec[1:], 15, ";"), true))
		common.InsertsFromArrays("ndcPharmClass", []string{
			"PharmClassName",
		}, common.OneToMultiArray(common.Derivatives(rec[1:], 16, ","), true))
		common.InsertsFromArrays("ndcDrugName", []string{
			"DrugName",
		}, common.OneToMultiArray(common.Derivatives(rec[1:], 3, ";"), true))
		return
	}

	fmt.Printf("DosageForm : %#v\n", common.Derivatives(rec[1:], 6, ";"))
	fmt.Printf("ProductType: %#v\n", common.Derivatives(rec[1:], 2, ";"))
	fmt.Printf("Route : %#v\n", common.Derivatives(rec[1:], 7, ";"))
	fmt.Printf("StrengthUnit : %#v\n", common.Derivatives(rec[1:], 15, ";"))
	fmt.Printf("PharmClasses : %#v\n", common.Derivatives(rec[1:], 16, ","))
	fmt.Printf("DrugName : %#v\n", common.Derivatives(rec[1:], 3, ";"))
	//fmt.Printf("%#v", rec[1:])
}

func newTsv(tsvFile string, rec [][]string, sourceField int, delim string) {
	fmt.Printf(" - Writing new %s\n", tsvFile)
	err := common.TsvFromArrays(
		"data/"+tsvFile,
		common.PrependUniqueIds(
			common.OneToMultiArray(common.Derivatives(rec[1:], sourceField, delim), true),
		),
	)
	if err != nil {
		panic(err)
	}
}

func mergeTsv(tsvFile string, rec [][]string, sourceField int, delim string) {
	fmt.Printf(" - Ingesting %s\n", tsvFile)
	raw, err := common.ReadTsv("data/" + tsvFile)
	if err != nil {
		panic(err)
	}
	newData := false
	newEntries := 0
	corpus := common.Derivatives(raw, 1, delim)
	updates := common.Derivatives(rec[1:], sourceField, delim)
	keys := common.CoerceSliceStringToInt(common.Derivatives(raw, 0, delim))
	maxVal := common.MaxIntSlice(keys)
	for _, u := range updates {
		uc := strings.TrimSpace(strings.ToUpper(u))
		if uc == "" {
			continue
		}

		pieces := strings.Split(uc, delim)
		for _, piece := range pieces {
			piece = strings.TrimSpace(piece)
			if !common.HasElement(corpus, piece) {
				newData = true
				maxVal++
				newEntries++
				if *Debug {
					fmt.Printf(" ! Found new element '%s' (%d)\n", piece, maxVal)
				}
				raw = append(raw, []string{fmt.Sprintf("%d", maxVal), piece})
				corpus = append(corpus, piece)
			}
		}
	}

	// Push out to TSV
	if newData {
		fmt.Printf(" - Writing updated %s with %d new entries\n", tsvFile, newEntries)
		err = common.TsvFromArrays(
			"data/"+tsvFile,
			raw,
		)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Printf(" - No new data to write to %s\n", tsvFile)
	}
}
