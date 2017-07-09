package controllers

import (
	"tournamentAPI/requests"
	"github.com/astaxie/beego"
	"tournamentAPI/services"
	"net/http"
	"encoding/json"
)

type TournamentController struct {
	beego.Controller
}

type emptyResponse struct {
}

type tournamentErrorResponse struct {
	TournamentId int
	Error string
}

func (this *TournamentController) AnnounceTournament() {

	var request = requests.NewAnnounceTournamentRequest()
	request.ExchangeMap(this.Input())
	errors := request.HasErrors()

	if errors != nil {
		this.formErrorResponse(errors[0], request.TournamentId, http.StatusBadRequest)
		return
	}

	service := services.NewTournamentService()
	err := service.AddTournament(request)

	if err != nil {
		this.formErrorResponse(err, request.TournamentId, http.StatusInternalServerError)
		return
	}

	this.Data["json"] = &changeBalanceResponse{}
	this.ServeJSON()
}

func (this *TournamentController) JoinTournament() {

	var request = requests.NewJoinTournamentRequest()
	request.ExchangeMap(this.Input())
	errors := request.HasErrors()

	if errors != nil {
		this.formErrorResponse(errors[0], request.TournamentId, http.StatusBadRequest)
		return
	}

	service := services.NewTournamentService()
	err := service.JoinTournament(request)

	if err != nil {
		this.formErrorResponse(err, request.TournamentId, http.StatusInternalServerError)
		return
	}

	this.Data["json"] = &changeBalanceResponse{}
	this.ServeJSON()
}

func (this *TournamentController) ResultTournament() {

	var request = requests.NewCloseTournamentRequest()
	jsonErr := json.Unmarshal(this.Ctx.Input.RequestBody, &request)
	if jsonErr != nil {
		this.formErrorResponse(jsonErr, request.TournamentId, http.StatusBadRequest)
		return
	}
	errors := request.HasErrors()

	if errors != nil {
		this.formErrorResponse(errors[0], request.TournamentId, http.StatusBadRequest)
		return
	}

	service := services.NewTournamentService()
	err := service.CloseTournament(request)

	if err != nil {
		this.formErrorResponse(err, request.TournamentId, http.StatusInternalServerError)
		return
	}

	this.Data["json"] = &changeBalanceResponse{}
	this.ServeJSON()
}


func (this *TournamentController) formErrorResponse(err error, tournamentId int, statusCode int) {
	this.Ctx.Output.SetStatus(statusCode)
	this.Data["json"] = &tournamentErrorResponse{
		TournamentId: tournamentId,
		Error:    err.Error(),
	}
	this.ServeJSON()
}