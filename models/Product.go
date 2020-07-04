package models

type Product struct {
	Id             int
	ChequeId       int
	Title          string
	SectionCode    string
	Price          string
	DiscountSum    int
	Discount       string
	DiscountStorno int
	ExtraCharge    string
	Number         string
	Sum            string
	MarkupSum      uint64
	MarkupStorno   bool
	Total          uint64
	IsDisPrice     bool
	IsDisDiscount  bool
	IsDisExCharge  bool
	IsDisNumber    bool
	IsStorno       bool
	//Cheque  *Cheque `orm:"rel(fk)"`       // RelForeignKey relation
}
