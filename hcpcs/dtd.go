package main

type HcpcsRecord struct {
	Code                             string     `fixed:"1-5"`
	ModiferCode                      string     `fixed:"4-5"`
	SequenceNumber                   string     `fixed:"6-10"`
	RecIdentificationCode            string     `fixed:"11-11"`
	LongDescription                  string     `fixed:"12-91"`
	ShortDescription                 string     `fixed:"92-119"`
	PricingIndicatorCode             string     `fixed:"120-121"`
	MultiplePricingIndicatorCode     string     `fixed:"128-128"`
	CoverageIssuesReferenceSection   string     `fixed:"129-134"`
	MedicareCarriersReferenceSection string     `fixed:"147-154"`
	StatuteNumber                    string     `fixed:"171-180"`
	LabCertificationCode             string     `fixed:"181-183"`
	CrossReferenceCode               string     `fixed:"205-209"`
	CoverageCode                     string     `fixed:"230-230"`
	AscPaymentGroupCode              string     `fixed:"231-232"`
	AscPaymentEffectiveDate          *HcpcsDate `fixed:"233-240"`
	MogPaymentGroupCode              string     `fixed:"241-243"`
	MogPaymentPolicyIndicator        string     `fixed:"244-244"`
	MogEffectiveDate                 *HcpcsDate `fixed:"245-252"`
	ProcessingNoteNumber             string     `fixed:"253-256"`
	BerensonEggersTosCode            string     `fixed:"257-259"`
	TosCode                          string     `fixed:"261-261"`
	AnesthesiaBaseUnitQuantity       int        `fixed:"266-268"`
	CodeAdded                        *HcpcsDate `fixed:"269-276"`
	EffectiveActionDate              *HcpcsDate `fixed:"277-284"`
	TerminationDate                  *HcpcsDate `fixed:"285-292"`
	ActionCode                       string     `fixed:"293-293"`
}

type HcpcsDate struct {
	Year  int `fixed:"1-4"`
	Month int `fixed:"5-6"`
	Day   int `fixed:"7-8"`
}
