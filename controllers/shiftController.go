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
			income = income + totalSum
			shift.Income = strconv.Itoa(income)
			shift.ShiftOpening = time.Now()
			shift.ShiftClosing = time.Now()
			o.Insert(&shift)
			handlers.SetTimer()

			kkmCash, _ := strconv.Atoi(kkm.Cash)
			kkm.Cash = strconv.Itoa(kkmCash + totalSum)
			o.Update(&kkm)

			cheque.TotalSum = chequeData.Cheque.TotalSum
			cheque.TotalDiscount = chequeData.Cheque.TotalDiscount
			cheque.TotalCharge = chequeData.Cheque.TotalCharge
			cheque.NDS = chequeData.Cheque.NDS
			cheque.PaymentType = chequeData.Cheque.PaymentType
			cheque.ChangeMoney = chequeData.Cheque.ChangeMoney
			cheque.Kkm_id, _ = strconv.Atoi(kkmId)
			cheque.Shift_id = shift.Id
			cheque.OperationType = "sale"
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
				o.Insert(&product)

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
		cheque.Kkm_id, _ = strconv.Atoi(kkmId)
		cheque.Shift_id = shifts[0].Id
		cheque.OperationType = "sale"
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
			o.Insert(&product)

		}

		s.Data["json"] = zreport
	}

	s.ServeJSON()
}

func (s *ShiftController) GetCheques() {
	kkmId := s.Ctx.Input.Param(":kkmId")
	o := orm.NewOrm()
	var cheques []models.Cheque
	o.QueryTable("cheque").Filter("kkm_id", kkmId).All(&cheques)
	//var products []models.Product
	//for i := 0; i < len(cheques); i++ {
	//	o.QueryTable("product").Filter("cheque_id", cheques[i].Id).All(&products)
	//	}
	//
	//var f interface{}
	//f = map[string]interface{}{
	//	"cheques":  cheques,
	//	"products": products,
	//}
	//s.Data["json"] = f
	s.Data["json"] = cheques
	s.ServeJSON()
}

func (s *ShiftController) PickCheque() {
	checkId := s.Ctx.Input.Param(":checkId")
	o := orm.NewOrm()
	var products []models.Product
	o.QueryTable("product").Filter("cheque_id", checkId).All(&products)
	s.Data["json"] = products
	s.ServeJSON()
}

func (s *ShiftController) ReturnSale() {
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
			shift.ShiftOpening = time.Now()
			shift.ShiftClosing = time.Now()
			o.Insert(&shift)
			handlers.SetTimer()

			kkmCash, _ := strconv.Atoi(kkm.Cash)
			kkm.Cash = strconv.Itoa(kkmCash - totalSum)
			o.Update(&kkm)

			cheque.TotalSum = chequeData.Cheque.TotalSum
			cheque.TotalDiscount = chequeData.Cheque.TotalDiscount
			cheque.TotalCharge = chequeData.Cheque.TotalCharge
			cheque.NDS = chequeData.Cheque.NDS
			cheque.PaymentType = chequeData.Cheque.PaymentType
			cheque.ChangeMoney = chequeData.Cheque.ChangeMoney
			cheque.Kkm_id, _ = strconv.Atoi(kkmId)
			cheque.Shift_id = shift.Id
			cheque.OperationType = "return"
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
		cheque.OperationType = "return"
		cheque.Kkm_id, _ = strconv.Atoi(kkmId)
		cheque.Shift_id = shifts[0].Id
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
			o.Insert(&product)

		}
		s.Data["json"] = lastZreport
	}

	s.ServeJSON()
}
