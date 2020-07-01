package models

type Nomenclature struct {
	Id              int
	OrganizationBin string
	Title           string
	Price           string
	Discount        string
	ExtraCharge     string
	Number          string
	Sum             string
	IsDisPrice      bool
	IsDisDiscount   bool
	IsDisExCharge   bool
	IsDisNumber     bool
	//Cheque  *Cheque `orm:"rel(fk)"`       // RelForeignKey relation
}
