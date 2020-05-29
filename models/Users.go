package models

import (
	"github.com/astaxie/beego/orm"
)

type Users struct {
	Id           int `orm:"auto"`
	FirstName    string
	LastName     string
	Email        string
	Password     string
	Role         string
	Organization []*Organization `orm:"reverse(many)"` // reverse relationship of fk
}

func Login(email, password string) bool {
	o := orm.NewOrm()
	var user Users
	o.QueryTable("users").Filter("email", email).All(&user)
	if user.Email == email && user.Password == password {
		return true
	}
	return false
}
