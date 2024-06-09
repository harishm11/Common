package models

import (
	"gorm.io/gorm"
)

type Bundle struct {
	gorm.Model     `swaggerignore:"true"`
	Transaction    Transaction        `json:"transaction" gorm:"foreignKey:BundleID"`
	Policy         Policy             `json:"policy" gorm:"foreignKey:BundleID"`
	CurrentCarrier CurrentCarrierInfo `json:"currentcarrier" gorm:"foreignKey:BundleID"`
	PolicyHolder   PolicyHolder       `json:"policyholder" gorm:"foreignKey:BundleID"`
	PolicyAddress  PolicyAddress      `json:"policyaddress" gorm:"foreignKey:BundleID"`
	Drivers        []Driver           `json:"driver" gorm:"foreignKey:BundleID"`
	Vehicles       []Vehicle          `json:"vehicle" gorm:"foreignKey:BundleID"`
}
