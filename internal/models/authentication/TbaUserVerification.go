package models

import (
	"time"
)

const (
	Administrator VerifyType = "administrator"
	Email         VerifyType = "email"
	SMS           VerifyType = "sms"
)

type VerifyType string

type TbaUserVerification struct {
	Id         uint       `gorm:"primaryKey;autoIncrement:true;"`
	UserId     uint       `gorm:"not null"`
	VerifyType VerifyType `gorm:"type:verify; not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time `gorm:"autoCreateTime:false"`
}

//func (ver *verifyType) Scan(value interface{}) error {
//	*ver = verifyType(value.([]byte))
//	return nil
//}
//func (ver verifyType) Value() (driver.Value, error) {
//	return string(ver), nil
//}

func (TbaUserVerification) TableName() string {
	return "tba_users_verification"
}
