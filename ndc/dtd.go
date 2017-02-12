package main

const (
	SqlPreamble = `
DROP TABLE IF EXISTS ndc;
CREATE TABLE ndc (
	  ProductID             VARCHAR(128) NOT NULL
	, ProductNDC            VARCHAR(128) NOT NULL
	, ProductTypeName       VARCHAR(128) NOT NULL
	, ProprietaryName       VARCHAR(128) NOT NULL
	, ProprietaryNameSuffix VARCHAR(128)
	, NonProprietaryName    VARCHAR(128)
	, DosageFormName        VARCHAR(128)
	, RouteName             VARCHAR(128)
	, StartMarketingDate    TIMESTAMP NULL
	, EndMarketingDate      TIMESTAMP NULL
	, MarketingCategoryName VARCHAR(128)
	, ApplicationNumber     VARCHAR(128)
	, LabelerName           VARCHAR(128)
	, SubstanceName         VARCHAR(128)
	, StrengthNumber        VARCHAR(100)
	, StrengthUnit          VARCHAR(128)
	, PharmClasses          VARCHAR(128)
	, DEASchedule           CHAR(5)

	, INDEX                ( ProductID )
	, INDEX                ( ProductNDC )
	, INDEX                ( ProductTypeName )
	, INDEX                ( ProprietaryName )
) ENGINE=InnoDB;
`
)

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
