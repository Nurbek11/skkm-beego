package models

type Kkm struct {
	Id              int
	OrganizationBin string
	Title           string
	Cash            string
	//Organization *Organization `orm:"rel(fk)"`   // RelForeignKey relation

}
