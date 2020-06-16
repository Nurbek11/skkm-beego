package models

import "time"

type Zreport struct {
	Id               int `orm:"auto"`
	DateTime   time.Time
	ShiftOpening     time.Time
	ShiftClosing     time.Time
	Shift            int
	StartSales       string
	StartPayouts     string
	StartSalesReturn string
	StartRefunds     string
	ShiftSales       string
	ShiftPayouts     string
	ShiftSalesReturn string
	ShiftRefunds     string
}
