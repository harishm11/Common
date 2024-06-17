package models

import "time"

type Ratebooks struct {
	StateCode                  string
	ProductCode                string
	RatebookCode               string
	RatebookActivationDate     time.Time
	RatebookNewBusinessEffDate time.Time
	RatebookRenBusinessEffDate time.Time
}
