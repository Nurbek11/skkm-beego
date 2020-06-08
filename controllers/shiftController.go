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
		shift.KkmId, _ = strconv.Atoi(kkmId)
		shift.Payouts = "0"
		shift.Depositing = "0"
		shift.Withdrawing = "0"
		shift.Income = "0"
		shift.ShiftOpening = time.Now()
		shift.ShiftClosing = time.Now()
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
		kkmCash, err := strconv.Atoi(kkm.Cash)
		kkm.Cash = strconv.Itoa(kkmCash + amount)
		newAmount, err := strconv.Atoi(shift[0].Depositing)
		if err == nil {
			newAmount = newAmount + amount
			shift[0].Depositing = strconv.Itoa(newAmount)
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
		kkmCash, _ := strconv.Atoi(kkm.Cash)
		if kkmCash < amount {
			s.Data["json"] = "insufficient funds"
		} else {
			kkm.Cash = strconv.Itoa(kkmCash - amount)
			newAmount, err := strconv.Atoi(shift[0].Withdrawing)
			if err == nil {
				newAmount = newAmount + amount
				shift[0].Withdrawing = strconv.Itoa(newAmount)
			}
			o.Update(&shift[0])
			o.Update(&kkm)
			s.Data["json"] = kkm
		}
	}
	s.ServeJSON()
}

func (s *ShiftController) ShowZreport() {
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
			"OverInfo": map[string]string{
				"address": 	organization.Address,
				"bin":     strconv.Itoa(organization.Bin),
				"shift_number":   strconv.Itoa(shift[0].Id),
				"cash":kkm.Cash,
				"creationTime":time.Now().String(),

			},
			"infoAtTheBeginningOfTheShift": map[string]string{
				"openingOfTheShift":shift[0].ShiftOpening.String(),
				"sales":"sales",
				"payouts":"payouts",
				"returnOfSales":"returnOfSales",
				"refunds":"refunds",

			},
			"info": map[string]string{

				"cash":"cash",
				"making":"making",
				"withdrawing":"withdrawing",
			},
			"infoForTheCurrentShift": map[string]string{
				"closingOfTheShift":shift[0].ShiftClosing.String(),
				"sales":"sales",
				"payouts":"payouts",
				"returnOfSales":"returnOfSales",
				"refunds":"refunds",
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
	var shifts []models.Shift
	var shift models.Shift
	var kkm models.Kkm
	o.QueryTable("kkm").Filter("id", kkmId).All(&kkm)
	o.QueryTable("shift").Filter("is_open_shift", true).All(&shifts)
	if len(shifts) == 0 {
		shift.IsOpenShift = true
		shift.KkmId, _ = strconv.Atoi(kkmId)
		shift.Payouts = "0"
		shift.Depositing = "0"
		shift.Withdrawing = "0"
		income, _ := strconv.Atoi(shift.Income)
		totalSum, _ := strconv.Atoi(chequeData.TotalSum)
		income = income+totalSum
		shift.Income = strconv.Itoa(income)
		shift.ShiftOpening = time.Now()
		shift.ShiftClosing = time.Now()
		o.Insert(&shift)
		handlers.SetTimer()

		kkmCash, _ := strconv.Atoi(kkm.Cash)
		kkm.Cash = strconv.Itoa(kkmCash+totalSum)
		o.Update(&kkm)

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

	} else  {
		income, err := strconv.Atoi(shifts[0].Income)
		if err == nil {
			totalSum, _ := strconv.Atoi(chequeData.TotalSum)
			income = income+totalSum
			shifts[0].Income = strconv.Itoa(income)
			kkmCash, _ := strconv.Atoi(kkm.Cash)
			kkm.Cash = strconv.Itoa(kkmCash+totalSum)
			o.Update(&kkm)
			o.Update(&shifts[0])
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
