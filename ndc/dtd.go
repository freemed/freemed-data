package main

const (
	SqlPreamble = `
DROP TABLE IF EXISTS ndc;
CREATE TABLE ndc (
	  ProductID             VARCHAR(128) NOT NULL
	, ProductNDC            VARCHAR(128) NOT NULL
	, ProductTypeName       VARCHAR(128) NOT NULL
	, ProprietaryName       VARCHAR(1024) NOT NULL
	, ProprietaryNameSuffix VARCHAR(128)
	, NonProprietaryName    TEXT
	, DosageFormName        VARCHAR(128)
	, RouteName             VARCHAR(1024)
	, StartMarketingDate    DATE NULL
	, EndMarketingDate      DATE NULL
	, MarketingCategoryName VARCHAR(128)
	, ApplicationNumber     VARCHAR(128)
	, LabelerName           VARCHAR(128) CHARACTER SET utf8
	, SubstanceName         TEXT
	, StrengthNumber        TEXT
	, StrengthUnit          TEXT
	, PharmClasses          TEXT
	, DEASchedule           ENUM ( '', 'CI', 'CII', 'CIII', 'CIV', 'CV' ) NOT NULL DEFAULT ''

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
