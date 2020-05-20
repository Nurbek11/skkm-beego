package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/skkm-beego/models"
)

// Operations about JWT
type JWTController struct {
	beego.Controller
}

// @Title CreateToken
// @Description create tokens
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.JWT
// @Failure 403 body is empty
// @router / [post]
func (j *JWTController) Post() {
	var user models.Users
	json.Unmarshal(j.Ctx.Input.RequestBody, &user)
	token := models.AddToken(user, j.Ctx.Input.Domain())
	j.Data["json"] = map[string]string{"token" :token}
	j.ServeJSON()

}


