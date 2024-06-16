package models

import "time"

type RateResponse struct {
	AccountNumber     int
	PolicyNumber      int
	TransactionNumber int
	Carrier           string
	Package           string
	Amount            float32
	EffDate           time.Time
	VehDetails        []VehData
	DrvCount          int
	DiscDetails       []DiscData
	RateStatus        string
	RateMessage       string
}
type CvgData struct {
	CoverageCode   string
	CoverageOption string
	CvgSymbol      string
	Amount         float32
}

type VehData struct {
	Vehid      uint
	VehYear    int
	VehMake    string
	VehModel   string
	Amount     float32
	CarSymbol  string
	CvgDetails []CvgData
}

type DiscData struct {
	DiscountCode string
}
