package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/dgrijalva/jwt-go"
	"github.com/skkm-beego/handlers"
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
		s.Data["json"] = shift
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
		shift[0].ShiftClosing = time.Now()
		o.Update(&shift[0])



		kkm.ShiftClosed = true
		o.Update(&kkm)

		elements := map[string]map[string]string{
			"OverInfo": {
				"address":           organization.Address,
				"bin":               organization.Bin,
				"shift_number":      strconv.Itoa(shift[0].Id),
				"cash":              kkm.Cash,
				"depositing":        shift[0].Depositing,
				"withdrawing":       shift[0].Withdrawing,
				"openingOfTheShift": shift[0].ShiftOpening,
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

func (s *MainController) ReturnXreport() {
	o := orm.NewOrm()
	orgBin := s.Ctx.Input.Param(":orgBin")
	kkmId := s.Ctx.Input.Param(":kkmId")
	var organization models.Organization
	o.QueryTable("organization").Filter("bin", orgBin).All(&organization)
	var kkm models.Kkm
	o.QueryTable("kkm").Filter("id", kkmId).All(&kkm)
	var shift []models.Shift
	o.QueryTable("shift").Filter("is_open_shift", true).All(&shift)

	elements := map[string]map[string]string{
		"OverInfo": {
			"address":           organization.Address,
			"bin":               organization.Bin,
			"shift_number":      strconv.Itoa(shift[0].Id),
			"cash":              kkm.Cash,
			"depositing":        shift[0].Depositing,
			"withdrawing":       shift[0].Withdrawing,
			"income":            shift[0].Income,
			"payouts":           shift[0].Payouts,
			"openingOfTheShift": shift[0].ShiftOpening,
			"closingOfTheShift": shift[0].ShiftClosing.String(),
		},
	}

	s.Data["json"] = elements
	s.ServeJSON()

}

func (s *MainController) GetProducts(){
	o := orm.NewOrm()
	orgBin := s.Ctx.Input.Param(":orgBin")
	var nomenclature []models.Nomenclature
	o.QueryTable("nomenclature").Filter("organization_bin", orgBin).All(&nomenclature)
	s.Data["json"] = nomenclature
	s.ServeJSON()
}

func (s *MainController) PickProduct(){
	o := orm.NewOrm()
	orgBin := s.Ctx.Input.Param(":orgBin")
	productId := s.Ctx.Input.Param(":productId")
	var nomenclature models.Nomenclature
	o.QueryTable("nomenclature").Filter("organization_bin", orgBin).Filter("id",productId).All(&nomenclature)
	s.Data["json"] = nomenclature
	s.ServeJSON()
}

func (s *MainController) EditProduct(){
	o := orm.NewOrm()
	orgBin := s.Ctx.Input.Param(":orgBin")
	productId := s.Ctx.Input.Param(":productId")
	var editedProduct models.Nomenclature
	json.Unmarshal([]byte(s.Ctx.Input.RequestBody), &editedProduct)
	var nomenclature models.Nomenclature
	o.QueryTable("nomenclature").Filter("organization_bin", orgBin).Filter("id",productId).All(&nomenclature)
	nomenclature.Title = editedProduct.Title
	nomenclature.Price = editedProduct.Price
	nomenclature.Discount = editedProduct.Discount
	nomenclature.ExtraCharge = editedProduct.ExtraCharge
	nomenclature.Sum = editedProduct.Sum
	nomenclature.IsDisPrice = editedProduct.IsDisPrice
	nomenclature.IsDisDiscount = editedProduct.IsDisDiscount
	nomenclature.IsDisNumber = editedProduct.IsDisNumber
	nomenclature.IsDisExCharge = editedProduct.IsDisExCharge
	nomenclature.QuantityInStock = editedProduct.QuantityInStock
	o.Update(&nomenclature)
	s.Data["json"] = nomenclature
	s.ServeJSON()
}

func (s *MainController) RemoveProduct(){
	o := orm.NewOrm()
	orgBin := s.Ctx.Input.Param(":orgBin")
	productId := s.Ctx.Input.Param(":productId")
	var nomenclature models.Nomenclature
	o.QueryTable("nomenclature").Filter("organization_bin", orgBin).Filter("id",productId).All(&nomenclature)
	o.Delete(&nomenclature)
	s.Data["json"] = productId
	s.ServeJSON()
}

func (s *MainController) CreateProduct(){
	o := orm.NewOrm()
	orgBin := s.Ctx.Input.Param(":orgBin")
	var requestProduct models.Nomenclature
	json.Unmarshal([]byte(s.Ctx.Input.RequestBody), &requestProduct)
	var nomenclature models.Nomenclature
	nomenclature.OrganizationBin = orgBin
	nomenclature.Title = requestProduct.Title
	nomenclature.Price = requestProduct.Price
	nomenclature.Discount = requestProduct.Discount
	nomenclature.ExtraCharge = requestProduct.ExtraCharge
	nomenclature.Sum = requestProduct.Sum
	nomenclature.IsDisPrice = requestProduct.IsDisPrice
	nomenclature.IsDisDiscount = requestProduct.IsDisDiscount
	nomenclature.IsDisNumber = requestProduct.IsDisNumber
	nomenclature.IsDisExCharge = requestProduct.IsDisExCharge
	nomenclature.QuantityInStock = requestProduct.QuantityInStock
	o.Insert(&nomenclature)
	s.Data["json"] = nomenclature
	s.ServeJSON()
}

func (s *MainController) ProbitCheque() {
	o := orm.NewOrm()
	kkmId := s.Ctx.Input.Param(":kkmId")
	orgBin := s.Ctx.Input.Param(":orgBin")
	var chequeData models.ChequeData
	json.Unmarshal([]byte(s.Ctx.Input.RequestBody), &chequeData)
	var cheque models.Cheque
	var shifts []models.Shift
	var shift models.Shift
	var kkm models.Kkm
	o.QueryTable("kkm").Filter("id", kkmId).All(&kkm)
	o.QueryTable("shift").Filter("is_open_shift", true).All(&shifts)
	password := string(chequeData.Password)
	var nomenclature []models.Nomenclature
	if len(shifts) == 0 {
		if password != kkm.Password {
			s.Data["json"] = "not correct"
		} else {
			shift.IsOpenShift = true
			shift.KkmId, _ = strconv.Atoi(kkmId)
			shift.Payouts = "0"
			shift.Depositing = "0"
			shift.Withdrawing = "0"
			income, _ := strconv.Atoi(shift.Income)
			totalSum, _ := strconv.Atoi(chequeData.Cheque.TotalSum)
			income = income + totalSum
			shift.Income = strconv.Itoa(income)
			shift.ShiftOpening = time.Now().String()
			o.Insert(&shift)
			handlers.SetTimer()



			kkmCash, _ := strconv.Atoi(kkm.Cash)
			kkm.Cash = strconv.Itoa(kkmCash + totalSum)
			kkm.ShiftClosed = false
			kkm.ShiftNumber = kkm.ShiftNumber+1
			o.Update(&kkm)

			cheque.TotalSum = chequeData.Cheque.TotalSum
			cheque.TotalDiscount = chequeData.Cheque.TotalDiscount
			cheque.TotalCharge = chequeData.Cheque.TotalCharge
			cheque.NDS = chequeData.Cheque.NDS
			cheque.PaymentType = chequeData.Cheque.PaymentType
			cheque.ChangeMoney = chequeData.Cheque.ChangeMoney
			cheque.Domain = chequeData.Cheque.Domain
			cheque.Markup = chequeData.Cheque.Markup
			cheque.TaxType = chequeData.Cheque.TaxType
			cheque.TaxationType = chequeData.Cheque.TaxationType
			cheque.TaxPercent = chequeData.Cheque.TaxPercent
			cheque.TaxSum = chequeData.Cheque.TaxSum
			cheque.TaxIsInTotalSum = chequeData.Cheque.TaxIsInTotalSum
			cheque.CustomerEmail = chequeData.Cheque.CustomerEmail
			cheque.CustomerPhone = chequeData.Cheque.CustomerPhone
			cheque.PosTransactionId = chequeData.Cheque.PosTransactionId
            cheque.PaymentCash = chequeData.Cheque.PaymentCash
            cheque.PaymentCard = chequeData.Cheque.PaymentCard
            cheque.PaymentCredit = chequeData.Cheque.PaymentCredit


			cheque.Kkm_id, _ = strconv.Atoi(kkmId)
			cheque.ShiftId = shift.Id
			cheque.OperationType = "sale"
			cheque.OperatorName = s.GetAuthUser().FirstName
			cheque.OperatorCode = s.GetAuthUser().Id
			cheque.DateTime = time.Now()
			o.Insert(&cheque)
			var zreports []models.Zreport
			o.QueryTable("zreport").All(&zreports)
			if len(zreports) == 0 {
				var zreport models.Zreport
				zreport.Id = shift.Id
				zreport.ShiftId = shift.Id
				zreport.Cash = kkm.Cash
				zreport.StartSales = chequeData.Cheque.TotalSum
				zreport.ShiftSales = chequeData.Cheque.TotalSum
				o.Insert(&zreport)
			} else {
				var zreportLast = zreports[len(zreports)-1]
				var zreport models.Zreport
				totalSum, _ := strconv.Atoi(chequeData.Cheque.TotalSum)
				var shiftSales, _ = strconv.Atoi(zreportLast.ShiftSales)
				zreport.Id = shift.Id
				zreport.ShiftId = shift.Id
				zreport.Cash = kkm.Cash
				zreport.StartSales = strconv.Itoa(shiftSales + totalSum)
				zreport.ShiftSales = strconv.Itoa(shiftSales + totalSum)
				o.Insert(&zreport)
			}

			for i := 0; i < len(chequeData.Cheque.Goods); i++ {
				var product models.Product
				product.ChequeId = cheque.Id
				product.Title = chequeData.Cheque.Goods[i].GoodTitle
				product.Price = chequeData.Cheque.Goods[i].GoodPrice
				product.Discount = chequeData.Cheque.Goods[i].GoodDiscount
				product.ExtraCharge = chequeData.Cheque.Goods[i].GoodExtraCharge
				product.Number = chequeData.Cheque.Goods[i].GoodNumber
				product.Sum = chequeData.Cheque.Goods[i].GoodSum
				product.IsDisPrice = chequeData.Cheque.Goods[i].IsDisPrice
				product.IsDisDiscount = chequeData.Cheque.Goods[i].IsDisDiscount
				product.IsDisExCharge = chequeData.Cheque.Goods[i].IsDisExCharge
				product.IsDisNumber = chequeData.Cheque.Goods[i].IsDisNumber
				product.DiscountSum = chequeData.Cheque.Goods[i].DiscountSum
				product.DiscountStorno = chequeData.Cheque.Goods[i].DiscountStorno
				product.MarkupSum = chequeData.Cheque.Goods[i].MarkupSum
				product.MarkupStorno = chequeData.Cheque.Goods[i].MarkupStorno
				product.IsStorno = chequeData.Cheque.Goods[i].IsStorno
				product.Total = chequeData.Cheque.Goods[i].Total
				product.SectionCode = chequeData.Cheque.Goods[i].SectionCode
				o.Insert(&product)

				o.QueryTable("nomenclature").Filter("organization_bin",orgBin).Filter("title", chequeData.Cheque.Goods[i].GoodTitle).All(&nomenclature)
				if len(nomenclature) == 0 {
					var nomen models.Nomenclature
					nomen.OrganizationBin = orgBin
					nomen.Title = chequeData.Cheque.Goods[i].GoodTitle
					nomen.Price = chequeData.Cheque.Goods[i].GoodPrice
					nomen.Discount = chequeData.Cheque.Goods[i].GoodDiscount
					nomen.ExtraCharge = chequeData.Cheque.Goods[i].GoodExtraCharge
					nomen.Sum = chequeData.Cheque.Goods[i].GoodSum
					nomen.IsDisPrice = chequeData.Cheque.Goods[i].IsDisPrice
					nomen.IsDisDiscount = chequeData.Cheque.Goods[i].IsDisDiscount
					nomen.IsDisExCharge = chequeData.Cheque.Goods[i].IsDisExCharge
					nomen.IsDisNumber = chequeData.Cheque.Goods[i].IsDisNumber

					o.Insert(&nomen)
				}


			}
			s.Data["json"] = "shift is opened"
		}
	} else {
		income, err := strconv.Atoi(shifts[0].Income)
		if err == nil {
			totalSum, _ := strconv.Atoi(chequeData.Cheque.TotalSum)
			income = income + totalSum
			shifts[0].Income = strconv.Itoa(income)
			kkmCash, _ := strconv.Atoi(kkm.Cash)
			kkm.Cash = strconv.Itoa(kkmCash + totalSum)
			o.Update(&kkm)
			o.Update(&shifts[0])
		}
		cheque.TotalSum = chequeData.Cheque.TotalSum
		cheque.TotalDiscount = chequeData.Cheque.TotalDiscount
		cheque.TotalCharge = chequeData.Cheque.TotalCharge
		cheque.NDS = chequeData.Cheque.NDS
		cheque.PaymentType = chequeData.Cheque.PaymentType
		cheque.ChangeMoney = chequeData.Cheque.ChangeMoney
		cheque.Domain = chequeData.Cheque.Domain
		cheque.Markup = chequeData.Cheque.Markup
		cheque.TaxType = chequeData.Cheque.TaxType
		cheque.TaxationType = chequeData.Cheque.TaxationType
		cheque.TaxPercent = chequeData.Cheque.TaxPercent
		cheque.TaxSum = chequeData.Cheque.TaxSum
		cheque.TaxIsInTotalSum = chequeData.Cheque.TaxIsInTotalSum
		cheque.CustomerEmail = chequeData.Cheque.CustomerEmail
		cheque.CustomerPhone = chequeData.Cheque.CustomerPhone
		cheque.PosTransactionId = chequeData.Cheque.PosTransactionId
		cheque.PaymentCash = chequeData.Cheque.PaymentCash
		cheque.PaymentCard = chequeData.Cheque.PaymentCard
		cheque.PaymentCredit = chequeData.Cheque.PaymentCredit
		cheque.Kkm_id, _ = strconv.Atoi(kkmId)
		cheque.ShiftId = shifts[0].Id
		cheque.OperationType = "sale"
		cheque.OperatorName = s.GetAuthUser().FirstName
		cheque.OperatorCode = s.GetAuthUser().Id
		cheque.DateTime = time.Now()

		o.Insert(&cheque)

		totalSum, _ := strconv.Atoi(chequeData.Cheque.TotalSum)
		var zreport models.Zreport
		o.QueryTable("zreport").Filter("id", shifts[0].Id).All(&zreport)
		zreport.Cash = kkm.Cash
		var shiftSales, _ = strconv.Atoi(zreport.ShiftSales)
		zreport.StartSales = strconv.Itoa(shiftSales + totalSum)
		zreport.ShiftSales = strconv.Itoa(shiftSales + totalSum)
		o.Update(&zreport)

		for i := 0; i < len(chequeData.Cheque.Goods); i++ {
			var product models.Product
			product.ChequeId = cheque.Id
			product.Title = chequeData.Cheque.Goods[i].GoodTitle
			product.Price = chequeData.Cheque.Goods[i].GoodPrice
			product.Discount = chequeData.Cheque.Goods[i].GoodDiscount
			product.ExtraCharge = chequeData.Cheque.Goods[i].GoodExtraCharge
			product.Number = chequeData.Cheque.Goods[i].GoodNumber
			product.Sum = chequeData.Cheque.Goods[i].GoodSum
			product.IsDisPrice = chequeData.Cheque.Goods[i].IsDisPrice
			product.IsDisDiscount = chequeData.Cheque.Goods[i].IsDisDiscount
			product.IsDisExCharge = chequeData.Cheque.Goods[i].IsDisExCharge
			product.IsDisNumber = chequeData.Cheque.Goods[i].IsDisNumber
			product.DiscountSum = chequeData.Cheque.Goods[i].DiscountSum
			product.DiscountStorno = chequeData.Cheque.Goods[i].DiscountStorno
			product.MarkupSum = chequeData.Cheque.Goods[i].MarkupSum
			product.MarkupStorno = chequeData.Cheque.Goods[i].MarkupStorno
			product.IsStorno = chequeData.Cheque.Goods[i].IsStorno
			product.Total = chequeData.Cheque.Goods[i].Total
			product.SectionCode = chequeData.Cheque.Goods[i].SectionCode
			o.Insert(&product)

			o.QueryTable("nomenclature").Filter("organization_bin",orgBin).Filter("title", chequeData.Cheque.Goods[i].GoodTitle).All(&nomenclature)
			if len(nomenclature) == 0 {
				var nomen models.Nomenclature
				nomen.OrganizationBin = orgBin
				nomen.Title = chequeData.Cheque.Goods[i].GoodTitle
				nomen.Price = chequeData.Cheque.Goods[i].GoodPrice
				nomen.Discount = chequeData.Cheque.Goods[i].GoodDiscount
				nomen.ExtraCharge = chequeData.Cheque.Goods[i].GoodExtraCharge
				nomen.Sum = chequeData.Cheque.Goods[i].GoodSum
				nomen.IsDisPrice = chequeData.Cheque.Goods[i].IsDisPrice
				nomen.IsDisDiscount = chequeData.Cheque.Goods[i].IsDisDiscount
				nomen.IsDisExCharge = chequeData.Cheque.Goods[i].IsDisExCharge
				nomen.IsDisNumber = chequeData.Cheque.Goods[i].IsDisNumber

				o.Insert(&nomen)
			}

		}

		s.Data["json"] = zreport
	}

	s.ServeJSON()
}

func (s *MainController) ReturnSale() {
	o := orm.NewOrm()
	kkmId := s.Ctx.Input.Param(":kkmId")
	var chequeData models.ChequeData
	json.Unmarshal([]byte(s.Ctx.Input.RequestBody), &chequeData)
	s.Data["json"] = chequeData.Cheque
	var cheque models.Cheque
	var shifts []models.Shift
	var shift models.Shift
	var kkm models.Kkm
	o.QueryTable("kkm").Filter("id", kkmId).All(&kkm)
	o.QueryTable("shift").Filter("is_open_shift", true).All(&shifts)
	password := string(chequeData.Password)

	if len(shifts) == 0 {
		if password != kkm.Password {
			s.Data["json"] = "not correct"
		} else {
			shift.IsOpenShift = true
			shift.KkmId, _ = strconv.Atoi(kkmId)
			shift.Payouts = "0"
			shift.Depositing = "0"
			shift.Withdrawing = "0"
			income, _ := strconv.Atoi(shift.Income)
			totalSum, _ := strconv.Atoi(chequeData.Cheque.TotalSum)
			income = income - totalSum
			shift.Income = strconv.Itoa(income)
			shift.ShiftOpening = time.Now().String()
			o.Insert(&shift)
			handlers.SetTimer()

			kkmCash, _ := strconv.Atoi(kkm.Cash)
			kkm.Cash = strconv.Itoa(kkmCash - totalSum)
			kkm.ShiftClosed = false
			kkm.ShiftNumber = kkm.ShiftNumber+1
			o.Update(&kkm)

			cheque.TotalSum = chequeData.Cheque.TotalSum
			cheque.TotalDiscount = chequeData.Cheque.TotalDiscount
			cheque.TotalCharge = chequeData.Cheque.TotalCharge
			cheque.NDS = chequeData.Cheque.NDS
			cheque.PaymentType = chequeData.Cheque.PaymentType
			cheque.ChangeMoney = chequeData.Cheque.ChangeMoney
			cheque.Domain = chequeData.Cheque.Domain
			cheque.Markup = chequeData.Cheque.Markup
			cheque.TaxType = chequeData.Cheque.TaxType
			cheque.TaxationType = chequeData.Cheque.TaxationType
			cheque.TaxPercent = chequeData.Cheque.TaxPercent
			cheque.TaxSum = chequeData.Cheque.TaxSum
			cheque.TaxIsInTotalSum = chequeData.Cheque.TaxIsInTotalSum
			cheque.CustomerEmail = chequeData.Cheque.CustomerEmail
			cheque.CustomerPhone = chequeData.Cheque.CustomerPhone
			cheque.PosTransactionId = chequeData.Cheque.PosTransactionId
			cheque.PaymentCash = chequeData.Cheque.PaymentCash
			cheque.PaymentCard = chequeData.Cheque.PaymentCard
			cheque.PaymentCredit = chequeData.Cheque.PaymentCredit
			cheque.Kkm_id, _ = strconv.Atoi(kkmId)
			cheque.ShiftId = shift.Id
			cheque.OperationType = "return"
			cheque.OperatorName = s.GetAuthUser().FirstName
			cheque.OperatorCode = s.GetAuthUser().Id
			cheque.DateTime = time.Now()
			o.Insert(&cheque)

			var zreports []models.Zreport
			o.QueryTable("zreport").All(&zreports)
			if len(zreports) == 0 {
				var zreport models.Zreport
				zreport.Id = shift.Id
				zreport.ShiftId = shift.Id
				zreport.Cash = kkm.Cash
				zreport.StartSalesReturn = chequeData.Cheque.TotalSum
				zreport.ShiftSalesReturn = chequeData.Cheque.TotalSum
				o.Insert(&zreport)
			} else {
				var zreportLast = zreports[len(zreports)-1]
				var zreport models.Zreport
				totalSum, _ := strconv.Atoi(chequeData.Cheque.TotalSum)
				var shiftSales, _ = strconv.Atoi(zreportLast.ShiftSalesReturn)
				zreport.Id = shift.Id
				zreport.ShiftId = shift.Id
				zreport.Cash = kkm.Cash
				zreport.StartSalesReturn = strconv.Itoa(shiftSales + totalSum)
				zreport.ShiftSalesReturn = strconv.Itoa(shiftSales + totalSum)
				o.Insert(&zreport)
			}

			for i := 0; i < len(chequeData.Cheque.Goods); i++ {
				var product models.Product
				product.ChequeId = cheque.Id
				product.Title = chequeData.Cheque.Goods[i].GoodTitle
				product.Price = chequeData.Cheque.Goods[i].GoodPrice
				product.Discount = chequeData.Cheque.Goods[i].GoodDiscount
				product.ExtraCharge = chequeData.Cheque.Goods[i].GoodExtraCharge
				product.Number = chequeData.Cheque.Goods[i].GoodNumber
				product.Sum = chequeData.Cheque.Goods[i].GoodSum
				product.IsDisPrice = chequeData.Cheque.Goods[i].IsDisPrice
				product.IsDisDiscount = chequeData.Cheque.Goods[i].IsDisDiscount
				product.IsDisExCharge = chequeData.Cheque.Goods[i].IsDisExCharge
				product.IsDisNumber = chequeData.Cheque.Goods[i].IsDisNumber
				product.DiscountSum = chequeData.Cheque.Goods[i].DiscountSum
				product.DiscountStorno = chequeData.Cheque.Goods[i].DiscountStorno
				product.MarkupSum = chequeData.Cheque.Goods[i].MarkupSum
				product.MarkupStorno = chequeData.Cheque.Goods[i].MarkupStorno
				product.IsStorno = chequeData.Cheque.Goods[i].IsStorno
				product.Total = chequeData.Cheque.Goods[i].Total
				product.SectionCode = chequeData.Cheque.Goods[i].SectionCode
				o.Insert(&product)

			}
			s.Data["json"] = "shift is opened"
		}
	} else {
		income, err := strconv.Atoi(shifts[0].Income)
		if err == nil {
			totalSum, _ := strconv.Atoi(chequeData.Cheque.TotalSum)
			income = income - totalSum
			shifts[0].Income = strconv.Itoa(income)
			kkmCash, _ := strconv.Atoi(kkm.Cash)
			kkm.Cash = strconv.Itoa(kkmCash - totalSum)
			o.Update(&kkm)
			o.Update(&shifts[0])
		}
		cheque.TotalSum = chequeData.Cheque.TotalSum
		cheque.TotalDiscount = chequeData.Cheque.TotalDiscount
		cheque.TotalCharge = chequeData.Cheque.TotalCharge
		cheque.NDS = chequeData.Cheque.NDS
		cheque.PaymentType = chequeData.Cheque.PaymentType
		cheque.ChangeMoney = chequeData.Cheque.ChangeMoney
		cheque.Domain = chequeData.Cheque.Domain
		cheque.Markup = chequeData.Cheque.Markup
		cheque.TaxType = chequeData.Cheque.TaxType
		cheque.TaxationType = chequeData.Cheque.TaxationType
		cheque.TaxPercent = chequeData.Cheque.TaxPercent
		cheque.TaxSum = chequeData.Cheque.TaxSum
		cheque.TaxIsInTotalSum = chequeData.Cheque.TaxIsInTotalSum
		cheque.CustomerEmail = chequeData.Cheque.CustomerEmail
		cheque.CustomerPhone = chequeData.Cheque.CustomerPhone
		cheque.PosTransactionId = chequeData.Cheque.PosTransactionId
		cheque.PaymentCash = chequeData.Cheque.PaymentCash
		cheque.PaymentCard = chequeData.Cheque.PaymentCard
		cheque.PaymentCredit = chequeData.Cheque.PaymentCredit
		cheque.OperationType = "return"
		cheque.Kkm_id, _ = strconv.Atoi(kkmId)
		cheque.ShiftId = shifts[0].Id
		cheque.OperatorName = s.GetAuthUser().FirstName
		cheque.OperatorCode = s.GetAuthUser().Id
		cheque.DateTime = time.Now()
		o.Insert(&cheque)

		totalSum, _ := strconv.Atoi(chequeData.Cheque.TotalSum)
		var zreport models.Zreport
		o.QueryTable("zreport").Filter("id", shifts[0].Id).All(&zreport)
		var lastZreport models.Zreport
		o.QueryTable("zreport").Filter("id", shifts[0].Id-1).All(&lastZreport)
		zreport.Cash = kkm.Cash
		if lastZreport.ShiftSalesReturn == "" {
			if zreport.StartSalesReturn == "" {
				zreport.StartSalesReturn = chequeData.Cheque.TotalSum
				zreport.ShiftSalesReturn = chequeData.Cheque.TotalSum
				o.Update(&zreport)
			} else {
				var salesReturn, _ = strconv.Atoi(zreport.StartSalesReturn)
				zreport.StartSalesReturn = strconv.Itoa(salesReturn + totalSum)
				zreport.ShiftSalesReturn = strconv.Itoa(salesReturn + totalSum)
				o.Update(&zreport)
			}
		} else {
			if zreport.StartSalesReturn == "" {
				var salesReturn, _ = strconv.Atoi(lastZreport.ShiftSalesReturn)
				zreport.StartSalesReturn = strconv.Itoa(salesReturn + totalSum)
				zreport.ShiftSalesReturn = strconv.Itoa(salesReturn + totalSum)
				o.Update(&zreport)
			} else {
				var salesReturn, _ = strconv.Atoi(zreport.StartSalesReturn)
				zreport.StartSalesReturn = strconv.Itoa(salesReturn + totalSum)
				zreport.ShiftSalesReturn = strconv.Itoa(salesReturn + totalSum)
				o.Update(&zreport)
			}

		}

		for i := 0; i < len(chequeData.Cheque.Goods); i++ {
			var product models.Product
			product.ChequeId = cheque.Id
			product.Title = chequeData.Cheque.Goods[i].GoodTitle
			product.Price = chequeData.Cheque.Goods[i].GoodPrice
			product.Discount = chequeData.Cheque.Goods[i].GoodDiscount
			product.ExtraCharge = chequeData.Cheque.Goods[i].GoodExtraCharge
			product.Number = chequeData.Cheque.Goods[i].GoodNumber
			product.Sum = chequeData.Cheque.Goods[i].GoodSum
			product.IsDisPrice = chequeData.Cheque.Goods[i].IsDisPrice
			product.IsDisDiscount = chequeData.Cheque.Goods[i].IsDisDiscount
			product.IsDisExCharge = chequeData.Cheque.Goods[i].IsDisExCharge
			product.IsDisNumber = chequeData.Cheque.Goods[i].IsDisNumber
			product.DiscountSum = chequeData.Cheque.Goods[i].DiscountSum
			product.DiscountStorno = chequeData.Cheque.Goods[i].DiscountStorno
			product.MarkupSum = chequeData.Cheque.Goods[i].MarkupSum
			product.MarkupStorno = chequeData.Cheque.Goods[i].MarkupStorno
			product.IsStorno = chequeData.Cheque.Goods[i].IsStorno
			product.Total = chequeData.Cheque.Goods[i].Total
			product.SectionCode = chequeData.Cheque.Goods[i].SectionCode
			o.Insert(&product)

		}
		s.Data["json"] = lastZreport
	}

	s.ServeJSON()
}





