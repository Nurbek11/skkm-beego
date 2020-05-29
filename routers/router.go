package routers

import (
	"github.com/astaxie/beego"
	"github.com/skkm-beego/controllers"
	"github.com/skkm-beego/handlers"
)

func init() {
	//beego.Router("/login", &controllers.MainController{},"post:Login")

	ns := beego.NewNamespace("v1",
		beego.NSRouter("/login", &controllers.MainController{}, "post:Login"),
		beego.NSNamespace("/org",
			beego.NSRouter("/:orgId", &controllers.MainController{}, "get:PickOrg"),
			beego.NSRouter("/:orgId/:kkmId", &controllers.MainController{}, "get:PickKkm"),
			beego.NSRouter("/:orgId/:kkmId/open", &controllers.ShiftController{}, "post:OpenShift"),
			beego.NSRouter("/:orgId/:kkmId/isopen", &controllers.ShiftController{}, "get:IsOpenShift"),
			beego.NSRouter("/:orgId/:kkmId/closeShift", &controllers.ShiftController{}, "post:CloseShift"),
			beego.NSRouter("/:orgId/:kkmId/probitCheck",&controllers.ShiftController{},"post:ProbitCheque"),
			beego.NSRouter("/:orgId/:kkmId/depositCash", &controllers.ShiftController{}, "post:DepositCash"),
			beego.NSRouter("/:orgId/:kkmId/withdrawCash", &controllers.ShiftController{}, "post:WithdrawCash"),
			beego.NSRouter("/:orgId/:kkmId/Xreport",&controllers.ShiftController{},"get:ShowXreport"),
			),
		beego.NSBefore(handlers.Jwt),
	)
	beego.AddNamespace(ns)
}


