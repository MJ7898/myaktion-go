package model

import "gorm.io/gorm"

type Donation struct {
	gorm.Model
	CampaignID uint
	Amount float64 `json:"amount" gorm:"notNull;check:amount >= 1.0"`
	DonorName string `json:"donorName" gorm:"notNull;size:40"`
	ReceiptRequested bool `json:"receiptRequested" gorm:"notNull"`
	Status Status `json:"status" gorm:"notNull;type:ENUM('TRANSFERRED','IN_PROCESS')"`
	Account Account `json:"account" gorm:"embedded;embeddedPrefix:account_"`
}

type Status string

const (
	TRANSFERRED Status = "TRANSFERRED"
	IN_PROCESS Status = "IN_PROCESS"
)

