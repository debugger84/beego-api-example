package requests

import "github.com/astaxie/beego/validation"

type AnnounceTournamentRequest struct {
	TournamentId int
	Deposit int
	BaseRequest
}

func NewAnnounceTournamentRequest() *AnnounceTournamentRequest {
	var req = new(AnnounceTournamentRequest)

	return req
}



func (this *AnnounceTournamentRequest) ExchangeMap(values map[string][]string) error {
	var err error
	err = this.exchangeFieldInt("tournamentId", values, &this.TournamentId)
	if err != nil {
		return err
	}
	err = this.exchangeFieldInt("deposit", values, &this.Deposit)
	if err != nil {
		return err
	}

	return nil
}

func (this *AnnounceTournamentRequest) HasErrors() []error {
	valid := validation.Validation{}
	valid.Required(this.TournamentId, "tournamentId")
	valid.Required(this.Deposit, "deposit")
	valid.Min(this.Deposit, 1, "deposit")

	return this.returnErrors(valid)
}