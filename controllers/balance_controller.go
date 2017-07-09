package controllers

import (
	"tournamentAPI/requests"
	"github.com/astaxie/beego"
	"tournamentAPI/models"
	"tournamentAPI/services"
	"net/http"
)

type BalanceController struct {
	beego.Controller
}

type changeBalanceResponse struct {
}

type balanceErrorResponse struct {
	PlayerId string
	Error string
}

func (this *BalanceController) Charge() {

	var request = requests.NewChangeBalanceRequest()
	request.ExchangeMap(this.Input())
	errors := request.HasErrors()

	if errors != nil {
		this.formErrorResponse(errors[0], request.PlayerId, http.StatusBadRequest)
		return
	}

	service := services.NewBalanceService()
	err := service.DecreasePlayerBalance(request)

	if err != nil {
		this.formErrorResponse(err, request.PlayerId, http.StatusInternalServerError)
		return
	}

	this.Data["json"] = &changeBalanceResponse{}
	this.ServeJSON()
}

func (this *BalanceController) Fund() {

	var request = requests.NewChangeBalanceRequest()
	request.ExchangeMap(this.Input())
	errors := request.HasErrors()

	if errors != nil {
		this.formErrorResponse(errors[0], request.PlayerId, http.StatusBadRequest)
		return
	}

	service := services.NewBalanceService()
	err := service.IncreasePlayerBalance(request)

	if err != nil {
		this.formErrorResponse(err, request.PlayerId, http.StatusInternalServerError)
		return
	}

	this.Data["json"] = &changeBalanceResponse{}
	this.ServeJSON()
}

func (this *BalanceController) Balance() {

	var request = requests.NewGetBalanceRequest()
	request.ExchangeMap(this.Input())
	errors := request.HasErrors()
	if errors != nil {
		this.formErrorResponse(errors[0], request.PlayerId, http.StatusBadRequest)
		return
	}

	service := services.NewBalanceService()
	balance, err := service.GetPlayerBalance(request.PlayerId)

	if err != nil {
		this.formErrorResponse(err, request.PlayerId, http.StatusInternalServerError)
		return
	}

	this.Data["json"] = &models.Player{
		PlayerId: request.PlayerId,
		Balance: balance,
	}
	this.ServeJSON()
}

func (this *BalanceController) formErrorResponse(err error, playerId string, statusCode int) {
	this.Ctx.Output.SetStatus(statusCode)
	this.Data["json"] = &balanceErrorResponse{
		PlayerId: playerId,
		Error: err.Error(),
	}
	this.ServeJSON()
}