package requests

import (
	"github.com/astaxie/beego/validation"
)

type CloseTournamentRequest struct {
	TournamentId int `json:",string"`
	Winners []*Winner
	BaseRequest
}

type Winner struct {
	PlayerId string
	Prize int
}

func NewCloseTournamentRequest() *CloseTournamentRequest {
	var req = new(CloseTournamentRequest)

	return req
}


func (this *CloseTournamentRequest) HasErrors() []error {
	valid := validation.Validation{}
	valid.Required(this.TournamentId, "tournamentId")
	valid.Required(this.Winners, "winners")

	return this.returnErrors(valid)
}