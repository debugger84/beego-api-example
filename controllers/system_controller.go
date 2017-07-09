package controllers

import (
	"github.com/astaxie/beego"
	"tournamentAPI/basic"
)

type SystemController struct {
	beego.Controller
}

func (this *SystemController) Reset() {
	client := basic.GetRedisClient()
	client.FlushAll()
}
