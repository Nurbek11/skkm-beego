package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/dgrijalva/jwt-go"
	"github.com/skkm-beego/models"
	"strconv"
	"time"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

func (c *MainController) Login() {
	email := c.GetString("email")
	password := c.GetString("password")
	o := orm.NewOrm()
	var user models.Users
	o.QueryTable("users").Filter("email", email).All(&user)
	if email == "" || password == "" {
		c.Data["json"] = "Заполните все поля"
	} else if email == user.Email && password == user.Password {
		json.Unmarshal(c.Ctx.Input.RequestBody, user)
		token := models.AddToken(user, c.Ctx.Input.Domain())
		elements := map[string]map[string]string{
			"userInfo": map[string]string{
				"id":        strconv.Itoa(user.Id),
				"FirstName": user.FirstName,
				"LastName":  user.LastName,
				"email":     user.Email,
				"role":      user.Role,
			},
			"token": map[string]string{
				"value": token,
			},
		}
		c.Data["json"] = elements
	} else if email != user.Email || password != user.Password {
		c.Data["json"] = "Неверное имя пользователя или пароль"
	}
	c.ServeJSON()
}

func (c *MainController) GetAuthUser() models.Users {
	claims := c.GetToken().Claims.(jwt.MapClaims)
	id := claims["sub"].(float64)
	o := orm.NewOrm()
	var user models.Users
	o.QueryTable("users").Filter("id", id).All(&user)
	return user
}

func (c *MainController) GetToken() *jwt.Token {
	var tokenString string = c.Ctx.Input.Header("token")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(beego.AppConfig.String("HMACKEY")), nil

	})

	if err != nil {
		c.Ctx.Output.SetStatus(403)
		var responseBody models.APIResponse = models.APIResponse{403, err.Error()}
		resBytes, err := json.Marshal(responseBody)
		c.Ctx.Output.Body(resBytes)
		if err != nil {
			panic(err)
		}
	}

	return token

}

func (c *MainController) GetOrgs() {
	o := orm.NewOrm()
	var orgs []models.Organization
	o.QueryTable("organization").All(&orgs)
	c.Data["json"] = &orgs
	c.ServeJSON()
}

func (c *MainController) PickOrg() {
	orgBin := c.Ctx.Input.Param(":orgBin")

	o := orm.NewOrm()
	var org []models.Organization
	var kkms []models.Kkm
	o.QueryTable("organization").Filter("bin", orgBin).All(&org)
	o.QueryTable("kkm").Filter("organization_bin", orgBin).All(&kkms)
	if len(org) == 0 {
		c.Data["json"] = "there is no organization with such an ID"
	} else {

		c.Data["json"] = org[0]

	}
	c.ServeJSON()
}

func (c *MainController) GetKkms() {
	orgBin := c.Ctx.Input.Param(":orgBin")
	o := orm.NewOrm()
	var kkms []models.Kkm
	o.QueryTable("kkm").Filter("organization_bin", orgBin).All(&kkms)
	if len(kkms) == 0 {
		c.Data["json"] = "There is no kkm"
	} else {
		c.Data["json"] = &kkms
	}
	c.ServeJSON()
}

func (c *MainController) PickKkm() {
	password := c.GetString("password")
	kkmId := c.Ctx.Input.Param(":kkmId")
	o := orm.NewOrm()
	var kkm []models.Kkm
	o.QueryTable("kkm").Filter("id", kkmId).All(&kkm)
	if len(kkm) == 0 {
		c.Data["json"] = "there is no kkm with such an ID"
	} else {
		if kkm[0].Password != password {
			c.Data["json"] = "not correct"
		} else {
			elements := map[string]map[string]string{
				"token": map[string]string{
					"value": c.GetToken().Raw,
				},
				"user": map[string]string{
					"id":       strconv.Itoa(c.GetAuthUser().Id),
					"username": c.GetAuthUser().Email,
				},
				"data": map[string]string{
					"id":    strconv.Itoa(kkm[0].Id),
					"title": kkm[0].Title,
				},
			}
			c.Data["json"] = elements
		}
	}
	c.ServeJSON()
}

func (s *MainController) CloseShift() {
	o := orm.NewOrm()
	orgBin := s.Ctx.Input.Param(":orgBin")
	kkmId := s.Ctx.Input.Param(":kkmId")
	var zreportall []models.Zreport
	o.QueryTable("zreport").All(&zreportall)
	var organization models.Organization
	o.QueryTable("organization").Filter("bin", orgBin).All(&organization)
	var kkm models.Kkm
	o.QueryTable("kkm").Filter("id", kkmId).All(&kkm)
	var shift []models.Shift
	o.QueryTable("shift").Filter("is_open_shift", true).All(&shift)
	if len(shift) < 1 {
		s.Data["json"] = "No open shift"
	} else {
		var zreport = zreportall[len(zreportall)-1]
		zreport.CashierId = s.GetAuthUser().Id
		zreport.OrganizationId, _ = strconv.Atoi(orgBin)
		zreport.ShiftId = shift[0].Id
		zreport.Cash = kkm.Cash
		zreport.TimeOfCreation = time.Now()
		var cheques []models.Cheque
		o.QueryTable("cheque").Filter("kkm_id", kkmId).Filter("shift_id", shift[0].Id).Filter("operation_type", "sale").All(&cheques)
		var chequesReturn []models.Cheque
		o.QueryTable("cheque").Filter("kkm_id", kkmId).Filter("shift_id", shift[0].Id).Filter("operation_type", "return").All(&chequesReturn)
		var total = 0
		var totalReturn = 0
		for i := 0; i < len(cheques); i++ {
			totalSum, _ := strconv.Atoi(cheques[i].TotalSum)
			total = total + totalSum
		}
		for i := 0; i < len(chequesReturn); i++ {
			totalReturnSum, _ := strconv.Atoi(chequesReturn[i].TotalSum)
			totalReturn = totalReturn + totalReturnSum
		}

		if len(zreportall) > 1 {
			var penultZreport = zreportall[len(zreportall)-2]

			if total == 0 {
				zreport.StartSales = penultZreport.ShiftSales
				zreport.ShiftSales = penultZreport.ShiftSales
				o.Update(&zreport)
			} else {
				var startSales, _ = strconv.Atoi(zreport.StartSales)
				zreport.StartSales = strconv.Itoa(startSales - total)
				o.Update(&zreport)
			}

			if totalReturn == 0 {
				zreport.StartSalesReturn = penultZreport.ShiftSalesReturn
				zreport.ShiftSalesReturn = penultZreport.ShiftSalesReturn
				o.Update(&zreport)
			} else {
				var startSalesReturn, _ = strconv.Atoi(zreport.StartSalesReturn)
				zreport.StartSalesReturn = strconv.Itoa(startSalesReturn - totalReturn)
				o.Update(&zreport)
			}
		} else {
			var startSales, _ = strconv.Atoi(zreport.StartSales)
			zreport.StartSales = strconv.Itoa(startSales - total)
			var startSalesReturn, _ = strconv.Atoi(zreport.StartSalesReturn)
			zreport.StartSalesReturn = strconv.Itoa(startSalesReturn - totalReturn)
			o.Update(&zreport)
		}

		shift[0].IsOpenShift = false
		o.Update(&shift[0])

		elements := map[string]map[string]string{
			"OverInfo": {
				"address":           organization.Address,
				"bin":               organization.Bin,
				"shift_number":      strconv.Itoa(shift[0].Id),
				"cash":              kkm.Cash,
				"depositing":        shift[0].Depositing,
				"withdrawing":       shift[0].Withdrawing,
				"openingOfTheShift": shift[0].ShiftOpening.String(),
				"closingOfTheShift": shift[0].ShiftClosing.String(),
			},

			"infoAtTheBeginningOfTheShift": {
				"sales":         zreport.StartSales,
				"returnOfSales": zreport.StartSalesReturn,
			},
			"infoForTheCurrentShift": {
				"sales":         zreport.ShiftSales,
				"returnOfSales": zreport.ShiftSalesReturn,
			},
		}
		s.Data["json"] = elements
	}
	s.ServeJSON()
}
