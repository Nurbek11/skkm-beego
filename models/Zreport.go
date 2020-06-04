package models

import "time"

type Zreport struct{
	Id           int       `orm:"auto"`
    Shift        int
	StartSales        string
	StartPayouts      string
	StartSalesReturn  string
	StartRefunds      string
    ShiftSales        string
	ShiftPayouts      string
	ShiftSalesReturn  string
	ShiftRefunds      string
	TimeOfCreation time.Time


}
