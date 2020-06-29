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

	StartSalesReturn string

	ShiftSales       string

	ShiftSalesReturn string

}
