package models

import (
	"time"
)

type Shift struct {
	Id           int `orm:"auto"`
	KkmId        int
	ShiftNumber    uint32
	DocumentNumber uint32

	Income       string
	Payouts      string
	Withdrawing  string
	Depositing   string
	IsOpenShift  bool
	ShiftOpening string
	ShiftClosing time.Time

}

//func IsShiftOpen(shift *Shift) bool {
//	return shift.IsOpenShift
//}
