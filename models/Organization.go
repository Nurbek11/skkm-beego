package models

type Organization struct {
	Id    int `orm:"auto"`
	Title string
	Bin   int
	Address   string
	User  *Users `orm:"rel(fk)"`       // RelForeignKey relation
	//Kkm   []*Kkm `orm:"reverse(many)"` // reverse relationship of fk
}



