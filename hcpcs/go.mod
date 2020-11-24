module github.com/freemed/freemed-data/hcpcs

go 1.15

replace (
	github.com/freemed/freemed-data => ../../freemed-data
	github.com/freemed/freemed-data/common => ../../freemed-data/common
)

require (
	github.com/freemed/freemed-data/common v0.0.0-00010101000000-000000000000
	github.com/jbuchbinder/gofixedfield v0.0.0-20201102175957-fbedb1ea9d63
	github.com/jbuchbinder/gosimplehttp v0.0.0-20170815145554-20db4d78d11f
)
