package models

type Cheque struct {
	Id            int    `orm:"auto"`
	TotalSum      string
	ChangeMoney   string
	TotalDiscount string
	TotalCharge   string
	NDS           string
	PaymentType   string
	OperationType string
	//Product []*Product `orm:"reverse(many)"` // reverse relationship of fk
}
