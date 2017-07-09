package requests

import "github.com/astaxie/beego/validation"

type JoinTournamentRequest struct {
	TournamentId int
	PlayerId string
	BackerIds []string
	BaseRequest
}

func NewJoinTournamentRequest() *JoinTournamentRequest {
	var req = new(JoinTournamentRequest)

	return req
}



func (this *JoinTournamentRequest) ExchangeMap(values map[string][]string) error {
	var err error
	err = this.exchangeFieldInt("tournamentId", values, &this.TournamentId)
	if err != nil {
		return err
	}
	err = this.exchangeFieldString("playerId", values, &this.PlayerId)
	if err != nil {
		return err
	}

	err = this.exchangeFieldArrayOfStrings("backerId", values, &this.BackerIds)
	if err != nil {
		return err
	}

	return nil
}

func (this *JoinTournamentRequest) HasErrors() []error {
	valid := validation.Validation{}
	valid.Required(this.TournamentId, "tournamentId")
	valid.Required(this.PlayerId, "playerId")

	return this.returnErrors(valid)
}