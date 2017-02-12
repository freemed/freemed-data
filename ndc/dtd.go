package main

type ProductRecord struct {
	ProductID             string
	ProductNDC            string
	ProductTypeName       string
	ProprietaryName       string
	ProprietaryNameSuffix string
	NonProprietaryName    string
	DosageFormName        string
	RouteName             string
	StartMarketingDate    *NdcDate
	EndMarketingDate      *NdcDate
	MarketingCategoryName string
	ApplicationNumber     string
	LabelerName           string
	SubstanceName         string
	StrengthNumber        string
	StrengthUnit          string
	PharmClasses          string
	DEASchedule           string
}

type NdcDate struct {
	Year  int `fixed:"1-4"`
	Month int `fixed:"5-6"`
	Day   int `fixed:"7-8"`
}
