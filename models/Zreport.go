package models

import "time"

type Zreport struct {
	Id               int `orm:"auto"`
	CashierId        int
	OrganizationId   int
	ShiftId          int
	Cash             string
	TimeOfCreation   time.Time
	StartSales       string
	StartPayouts     string
	StartSalesReturn string
	StartRefunds     string
	ShiftSales       string
	ShiftPayouts     string
	ShiftSalesReturn string
	ShiftRefunds     string
}
