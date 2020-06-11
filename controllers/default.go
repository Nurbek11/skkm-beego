package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/dgrijalva/jwt-go"
	"github.com/skkm-beego/models"
	"strconv"
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

	kkmId := c.Ctx.Input.Param(":kkmId")
	o := orm.NewOrm()
	var kkm []models.Kkm
	o.QueryTable("kkm").Filter("id", kkmId).All(&kkm)
	if len(kkm) == 0 {
		c.Data["json"] = "there is no kkm with such an ID"
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
	c.ServeJSON()
}
