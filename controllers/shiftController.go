package controllers

import (
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
		shift.ShiftOpening = time.Now().String()
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



