package models

import (
	"time"
)

type Shift struct {
	Id           int `orm:"auto"`
	KkmId        int
	Income       string
	Payouts      string
	Withdrawing  string
	Depositing   string
	IsOpenShift  bool
	ShiftOpening time.Time
	ShiftClosing time.Time
	//`orm:"auto_now_add"`
}

//func IsShiftOpen(shift *Shift) bool {
//	return shift.IsOpenShift
//}
