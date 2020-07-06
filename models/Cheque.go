package models

import "time"

type Cheque struct {
	Id      int `orm:"auto"`
	Kkm_id  int
	ShiftId int

	PosTransactionId   string

	PaymentType        string
	OperationType      string

	DateTime           time.Time
	Domain             string
	MoneyPlacementType string
	TaxType            string
	TaxationType       string
	TaxPercent         uint32
	TaxSum             uint64
	TaxIsInTotalSum    bool

	OperatorName string
	OperatorCode int

	PaymentCash   int
	PaymentCard   int
	PaymentCredit int

	TotalSum      string

	ChangeMoney   string
	Markup        uint64
	TotalDiscount string
	TotalCharge   string

	CustomerEmail string
	CustomerPhone string
	TicketNumber  string

	OfflineTicketNumber string
	PrintedTicketNumber string

	ShiftDocumentNumber uint32

	QrCode     string
	IsCanceled bool

	NDS string

	//Product []*Product `orm:"reverse(many)"` // reverse relationship of fk
}
