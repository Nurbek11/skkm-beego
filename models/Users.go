package models

import (
	"github.com/astaxie/beego/orm"
)

type Users struct {
	Id       int `orm:"auto"`
	Username string
	Password string
	Role     string
	Organization []*Organization `orm:"reverse(many)"`   // reverse relationship of fk
}

func Login(username, password string) bool {
	o := orm.NewOrm()
	var user Users
	o.QueryTable("users").Filter("username",username).All(&user)
	if user.Username == username && user.Password == password {
		return true
	}
	return false
}


