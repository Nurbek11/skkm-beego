package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
	"github.com/skkm-beego/models"
	_ "github.com/skkm-beego/routers"
)

func init() {
	username := beego.AppConfig.String("postgresuser")
	password := beego.AppConfig.String("postgrespass")
	host := beego.AppConfig.String("postgreshost")
	dbname := beego.AppConfig.String("postgresdb")

	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", "user= "+username+" password= "+password+" host="+host+" port=5432 dbname= "+dbname+" sslmode=disable")
	orm.RegisterModel(new(models.Users), new(models.Organization), new(models.Kkm), new(models.Shift), new(models.Cheque), new(models.Product), new(models.Zreport), new(models.Nomenclature))
}

func main() {

	beego.Run()
}
