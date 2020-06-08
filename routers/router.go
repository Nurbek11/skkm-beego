package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/skkm-beego/controllers"
	"github.com/skkm-beego/handlers"
)

func init() {


	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type","token"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

	ns := beego.NewNamespace("v1",
		beego.NSRouter("/login", &controllers.MainController{}, "post:Login"),
		beego.NSNamespace("/org",
			beego.NSRouter("/",&controllers.MainController{},"get:GetOrgs"),
			beego.NSRouter("/:orgId", &controllers.MainController{}, "get:PickOrg"),
			beego.NSRouter("/:orgId/kkms", &controllers.MainController{}, "get:GetKkms"),
			beego.NSRouter("/:orgId/:kkmId", &controllers.MainController{}, "get:PickKkm"),
			beego.NSRouter("/:orgId/:kkmId/open", &controllers.ShiftController{}, "post:OpenShift"),
			beego.NSRouter("/:orgId/:kkmId/isopen", &controllers.ShiftController{}, "get:IsOpenShift"),
			beego.NSRouter("/:orgId/:kkmId/closeShift", &controllers.ShiftController{}, "post:CloseShift"),
			beego.NSRouter("/:orgId/:kkmId/probitCheck",&controllers.ShiftController{},"post:ProbitCheque"),
			beego.NSRouter("/:orgId/:kkmId/depositCash", &controllers.ShiftController{}, "post:DepositCash"),
			beego.NSRouter("/:orgId/:kkmId/withdrawCash", &controllers.ShiftController{}, "post:WithdrawCash"),
			beego.NSRouter("/:orgId/:kkmId/Zreport",&controllers.ShiftController{},"get:ShowZreport"),
			),
		beego.NSBefore(handlers.Jwt),
	)
	beego.AddNamespace(ns)
}


