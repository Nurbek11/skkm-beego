package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
	"github.com/skkm-beego/models"
	_ "github.com/skkm-beego/routers"
)

func init() {
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", "user=postgres password=nurbek32335 host=127.0.0.1 port=5432 dbname=skkm sslmode=disable")
	orm.RegisterModel(new(models.Users), new(models.Organization), new(models.Kkm), new(models.Shift), new(models.Cheque), new(models.Product))
}

func main() {

	beego.Run()
}
