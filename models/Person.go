package models

import (
	"time"
)

type Person struct {
	ID            uint      `gorm:"primaryKey;uniqueIndex;"`
	FirstName     string    `json:"FirstName" display:"First Name"`
	LastName      string    `json:"LastName" display:"Last Name"`
	MiddleName    string    `json:"MiddleName" display:"Middle Name"`
	DOB           time.Time `json:"DOB" display:"Date of Birth"`
	Gender        string    `json:"Gender" display:"Gender"`
	MaritalStatus string    `json:"MaritalStatus" display:"Marital Status"`
}
