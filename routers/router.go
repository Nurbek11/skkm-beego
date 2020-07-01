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
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type", "token"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

	ns := beego.NewNamespace("v1",
		beego.NSRouter("/login", &controllers.MainController{}, "post:Login"),
		beego.NSNamespace("/org",
			beego.NSRouter("/", &controllers.MainController{}, "get:GetOrgs"),
			beego.NSRouter("/:orgBin", &controllers.MainController{}, "get:PickOrg"),
			beego.NSRouter("/:orgBin/kkms", &controllers.MainController{}, "get:GetKkms"),
			beego.NSRouter("/:orgBin/kkms/:kkmId", &controllers.MainController{}, "post:PickKkm"),
			beego.NSRouter("/:orgBin/kkms/:kkmId/open", &controllers.ShiftController{}, "post:OpenShift"),
			beego.NSRouter("/:orgBin/kkms/:kkmId/isopen", &controllers.ShiftController{}, "get:IsOpenShift"),
			beego.NSRouter("/:orgBin/kkms/:kkmId/Xreport", &controllers.MainController{}, "get:ReturnXreport"),
			beego.NSRouter("/:orgBin/kkms/:kkmId/closeShift", &controllers.MainController{}, "post:CloseShift"),
			beego.NSRouter("/:orgBin/kkms/:kkmId/probitCheck", &controllers.ShiftController{}, "post:ProbitCheque"),
			beego.NSRouter("/:orgBin/kkms/:kkmId/returnSale", &controllers.ShiftController{}, "post:ReturnSale"),
			beego.NSRouter("/:orgBin/kkms/:kkmId/depositCash", &controllers.ShiftController{}, "post:DepositCash"),
			beego.NSRouter("/:orgBin/kkms/:kkmId/withdrawCash", &controllers.ShiftController{}, "post:WithdrawCash"),
			beego.NSRouter("/:orgBin/kkms/:kkmId/checks", &controllers.ShiftController{}, "get:GetCheques"),
			beego.NSRouter("/:orgBin/kkms/:kkmId/checks/:checkId", &controllers.ShiftController{}, "get:PickCheque"),
			beego.NSRouter("/:orgBin/kkms/:kkmId/products/", &controllers.MainController{}, "get:GetProducts"),
			beego.NSRouter("/:orgBin/kkms/:kkmId/products/:productId", &controllers.MainController{}, "get:PickProduct"),
			beego.NSRouter("/:orgBin/kkms/:kkmId/products/:productId/edit", &controllers.MainController{}, "post:EditProduct"),
		),
		beego.NSBefore(handlers.Jwt),
	)
	beego.AddNamespace(ns)
}
