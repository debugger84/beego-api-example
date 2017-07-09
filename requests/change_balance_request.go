package requests

import "github.com/astaxie/beego/validation"

type ChangeBalanceRequest struct {

	/**
	 * Identifier of the user, whose balance will be changed
	 */
	PlayerId string `form:"playerId"`

	/**
	 * Amount of points to change balance value
	 */
	Points int `form:"points"`

	BaseRequest
}

func NewChangeBalanceRequest() *ChangeBalanceRequest {
	var req = new(ChangeBalanceRequest)

	return req
}



func (this *ChangeBalanceRequest) ExchangeMap(values map[string][]string) error {
	var err error
	err = this.exchangeFieldString("playerId", values, &this.PlayerId)
	if err != nil {
		return err
	}
	err = this.exchangeFieldInt("points", values, &this.Points)
	if err != nil {
		return err
	}

	return nil
}

func (this *ChangeBalanceRequest) HasErrors() []error {
	valid := validation.Validation{}
	valid.Required(this.PlayerId, "playerId")
	valid.Required(this.Points, "points")
	valid.Min(this.Points, 1, "points")

	return this.returnErrors(valid)
}