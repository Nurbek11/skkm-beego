package models

type Kkm struct {
	Id       int
	Title    string
    Cash     int
	Organization *Organization `orm:"rel(fk)"`   // RelForeignKey relation

}
