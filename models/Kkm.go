package models

import (
	"time"
)

type Kkm struct {
	Id              int
	OfdId           uint32
	OfdReqNum       uint16
	OfdToken        uint32
	ShiftNumber     uint32
	ShiftOpenDate   time.Time
	ShiftClosed     bool
	PrintedNumber   uint64
	SerialNumber    string
	FnsKkmId        string
	Bin             string
	Title           string
	Address         string
	Password        string
	OrganizationBin string
	Cash            string
	//Organization *Organization `orm:"rel(fk)"`   // RelForeignKey relation
}
