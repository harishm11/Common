package models

import "time"

type RateFactors struct {
	StateCode              string
	ProductCode            string
	CoverageCode           string
	RatingItemCode         string
	RateVar1Code           string
	RateVar1Value          string
	RateVar2Code           string
	RateVar2Value          string
	RateVar3Code           string
	RateVar3Value          string
	RateVar4Code           string
	RateVar4Value          string
	RateVar5Code           string
	RateVar5Value          string
	RateVar6Code           string
	RateVar6Value          string
	RateVar7Code           string
	RateVar7Value          string
	RateVar8Code           string
	RateVar8Value          string
	RateFactor             float32
	RateActivationDate     time.Time
	RateNewBusinessEffDate time.Time
	RateRenBusinessEffDate time.Time
	RetiredRow             bool
}
