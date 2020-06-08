package models

type Kkm struct {
	Id             int
	OrganizationId int
	Title          string
	Cash           string
	//Organization *Organization `orm:"rel(fk)"`   // RelForeignKey relation

}
