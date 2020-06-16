package models

type Product struct {
	Id            int
	ChequeId      int
	Title         string
	Price         string
	Discount      string
	ExtraCharge   string
	Number        string
	Sum           string
	IsDisPrice    bool
	IsDisDiscount bool
	IsDisExCharge bool
	IsDisNumber   bool
	//Cheque  *Cheque `orm:"rel(fk)"`       // RelForeignKey relation
}
