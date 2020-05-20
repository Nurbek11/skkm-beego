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
			o.Update(&kkm)
			s.Data["json"] = kkm
		}
	}
	s.ServeJSON()
}

func (s *ShiftController) ShowXreport() {
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
	var chequeData models.ChequeData
	json.Unmarshal([]byte(s.Ctx.Input.RequestBody), &chequeData)
	s.Data["json"] = chequeData
	s.ServeJSON()
}
