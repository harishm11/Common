package models

type Address struct {
	ID                    uint   `gorm:"primaryKey;uniqueIndex;"`
	AddressLine1          string `json:"AddressLine1"`
	AddressLine2          string `json:"AddressLine2"`
	City                  string `json:"City"`
	County                string `json:"County"`
	State                 string `json:"State"`
	ZipCode               string `json:"ZipCode"`
	StandardizedIndicator bool
}
