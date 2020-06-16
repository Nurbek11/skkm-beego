package models

type Organization struct {
	Bin      string `orm:"pk"`

	Title    string
	Address  string
	//User  *Users `orm:"rel(fk)"`       // RelForeignKey relation
	//Kkm   []*Kkm `orm:"reverse(many)"` // reverse relationship of fk
}
