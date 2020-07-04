package handlers

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/skkm-beego/models"
	"time"
)

const timeInHours = 24

var timer time.Timer

func SetTimer() () {
	timer := time.NewTimer(time.Second * timeInHours)
	go func() {
		<-timer.C
		fmt.Println("You have to refresh the shift")
		o := orm.NewOrm()
		var shift []models.Shift
		o.QueryTable("shift").Filter("is_open_shift", true).All(&shift)
		if len(shift) != 0 {
			shift[0].IsOpenShift = false
			shift[0].ShiftClosing = time.Now()
			o.Update(&shift[0])
		}
	}()


}
