package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/skkm-beego/handlers"
	"github.com/skkm-beego/models"
	"strconv"
	"time"
)

type ShiftController struct {
	beego.Controller
}

func (s *ShiftController) OpenShift() {
	kkmId := s.Ctx.Input.Param(":kkmId")
	o := orm.NewOrm()
	var shifts []models.Shift
	var shift models.Shift
	o.QueryTable("shift").Filter("is_open_shift", true).All(&shifts)

	if len(shifts) == 0 {
		shift.IsOpenShift = true
		shift.KkmId, _ = s.GetInt(kkmId)
		shift.Making = "0"
		shift.Payouts = "0"
		shift.Income = "0"
		o.Insert(&shift)
		s.Data["json"] = &shift
		handlers.SetTimer()
	} else {
		s.Data["json"] = "is already opened"
	}
	s.ServeJSON()
}

func (s *ShiftController) IsOpenShift() {
	o := orm.NewOrm()
	var shift []models.Shift
	o.QueryTable("shift").Filter("is_open_shift", true).All(&shift)
	if len(shift) == 0 {
		s.Data["json"] = "No open shift"
	} else {
		s.Data["json"] = &shift
	}
	s.ServeJSON()
}

func (s *ShiftController) CloseShift() {
	o := orm.NewOrm()
	var shift []models.Shift
	o.QueryTable("shift").Filter("is_open_shift", true).All(&shift)
	if len(shift) == 0 {
		s.Data["json"] = "No open shift"
	} else {
		shift[0].IsOpenShift = false
		shift[0].ShiftClosing = time.Now()
		o.Update(&shift[0])
		s.Data["json"] = &shift[0]
	}
	s.ServeJSON()
}

func (s *ShiftController) DepositCash() {
	amount, _ := s.GetInt("amount")
	kkmId := s.Ctx.Input.Param(":kkmId")
	o := orm.NewOrm()
	var shift []models.Shift
	o.QueryTable("shift").Filter("is_open_shift", true).All(&shift)
	if len(shift) == 0 {
		s.Data["json"] = "Shift is closed"
	} else {
		var kkm models.Kkm
		o.QueryTable("kkm").Filter("id", kkmId).All(&kkm)
		kkm.Cash = kkm.Cash + amount
		newAmount, err := strconv.Atoi(shift[0].Making)
		if err == nil {
			newAmount = newAmount + amount
			shift[0].Making = strconv.Itoa(newAmount)
		}
		o.Update(&shift[0])
		o.Update(&kkm)
		s.Data["json"] = kkm
	}
	s.ServeJSON()

}

func (s *ShiftController) WithdrawCash() {
	amount, _ := s.GetInt("amount")
	kkmId := s.Ctx.Input.Param(":kkmId")
	o := orm.NewOrm()
	var shift []models.Shift
	o.QueryTable("shift").Filter("is_open_shift", true).All(&shift)
	if len(shift) == 0 {
		s.Data["json"] = "Shift is closed"
	} else {
		var kkm models.Kkm
		o.QueryTable("kkm").Filter("id", kkmId).All(&kkm)
		if kkm.Cash < amount {
			s.Data["json"] = "insufficient funds"
		} else {
			kkm.Cash = kkm.Cash - amount
			newAmount, err := strconv.Atoi(shift[0].Payouts)
			if err == nil {
				newAmount = newAmount + amount
				shift[0].Payouts = strconv.Itoa(newAmount)
			}
			o.Update(&shift[0])
			o.Update(&kkm)
			s.Data["json"] = kkm
		}
	}
	s.ServeJSON()
}

func (s *ShiftController) ShowZreport(){
	o := orm.NewOrm()
	orgId := s.Ctx.Input.Param(":orgId")
	kkmId := s.Ctx.Input.Param(":kkmId")
	var organization models.Organization
	o.QueryTable("organization").Filter("id", orgId).All(&organization)
	var kkm models.Kkm
	o.QueryTable("kkm").Filter("id", kkmId).All(&kkm)
	var shift []models.Shift
	o.QueryTable("shift").Filter("is_open_shift", true).All(&shift)
	if len(shift) == 0 {
		s.Data["json"] = "No open shift,please open it"
	} else {
		elements := map[string]map[string]string{
			"organizationInfo": map[string]string{
				"address": organization.Address,
				"bin":     strconv.Itoa(organization.Bin),
				"shift":   strconv.Itoa(shift[0].Id),
			},
			"salesInfo": map[string]string{
				"sales":             "nil",
				"salesReturn":       "nil",
				"salesReceipts":     "nil",
				"sale":              "nil",
				"salesReturnChecks": "nil",
				"returnOfSale":      "nil",
				"payouts":           "nil",
				"making":            "nil",
				"cash":              strconv.Itoa(kkm.Cash),
				"income":            "nil",
			},
		}
		s.Data["json"] = elements
	}
	s.ServeJSON()

}

func (s *ShiftController) ProbitCheque() {
	o := orm.NewOrm()
	kkmId := s.Ctx.Input.Param(":kkmId")
	var chequeData models.ChequeData
	json.Unmarshal([]byte(s.Ctx.Input.RequestBody), &chequeData)
	var cheque models.Cheque
	var shift []models.Shift
	var kkm models.Kkm
	o.QueryTable("kkm").Filter("id", kkmId).All(&kkm)
	o.QueryTable("shift").Filter("is_open_shift", true).All(&shift)
	if len(shift) == 0 {
		s.Data["json"] = "Shift is closed,please open it"
	} else {
		income, err := strconv.Atoi(shift[0].Income)
		if err == nil {
			totalSum, _ := strconv.Atoi(chequeData.TotalSum)
			income = income+totalSum
			shift[0].Income = strconv.Itoa(income)
			kkm.Cash = kkm.Cash+totalSum
			o.Update(&kkm)
			o.Update(&shift[0])
		}
		cheque.TotalSum = chequeData.TotalSum
		cheque.TotalDiscount = chequeData.TotalDiscount
		cheque.TotalCharge = chequeData.TotalCharge
		cheque.NDS = chequeData.NDS
		cheque.PaymentType = chequeData.PaymentType
		cheque.ChangeMoney = chequeData.ChangeMoney
		o.Insert(&cheque)

		for i := 0; i < len(chequeData.Goods); i++ {
			var product models.Product
			product.ChequeId = cheque.Id
			product.Title = chequeData.Goods[i].GoodTitle
			product.Price = chequeData.Goods[i].GoodPrice
			product.Discount = chequeData.Goods[i].GoodDiscount
			product.ExtraCharge = chequeData.Goods[i].GoodExtraCharge
			product.Number = chequeData.Goods[i].GoodNumber
			product.Sum = chequeData.Goods[i].GoodSum
			product.IsDisPrice = chequeData.Goods[i].IsDisPrice
			product.IsDisDiscount = chequeData.Goods[i].IsDisDiscount
			product.IsDisExCharge = chequeData.Goods[i].IsDisExCharge
			product.IsDisNumber = chequeData.Goods[i].IsDisNumber
			o.Insert(&product)

		}
		s.Data["json"] = cheque
	}

	s.ServeJSON()
}
