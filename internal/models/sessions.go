package models

import "time"

type Session struct {
	Token  string    `gorm:"type:char(43);primaryKey"`
	Data   []byte    `gorm:"type:blob;not null"`
	Expiry time.Time `gorm:"type:timestamp(6);not null;index"`
}
