package models

type Product struct {
	GoodId            int `orm:"pk"`
	ChequeId      int
	GoodTitle         string
	GoodPrice         string
	GoodDiscount      string
	GoodExtraCharge   string
	GoodNumber        string
	GoodSum           string
	IsDisPrice    bool
	IsDisDiscount bool
	IsDisExCharge bool
	IsDisNumber   bool
	//Cheque  *Cheque `orm:"rel(fk)"`       // RelForeignKey relation
}
