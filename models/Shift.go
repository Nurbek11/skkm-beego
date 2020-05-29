package models

import
(
	"time"
)

type Shift struct {
	Id           int       `orm:"auto"`
	KkmId        int
	Income       string
	Payouts      string
	Making       string
	IsOpenShift  bool
	ShiftOpening time.Time `orm:"auto_now_add"`
	ShiftClosing time.Time `orm:"auto_now_add"`
	//Kkm *Kkm `orm:"rel(fk)"`
}


func IsShiftOpen(shift *Shift) bool {
	return shift.IsOpenShift
}