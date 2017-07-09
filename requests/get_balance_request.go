package requests

import (
	"github.com/astaxie/beego/validation"
)

type GetBalanceRequest struct {

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

func NewGetBalanceRequest() *GetBalanceRequest {
	var req = new(GetBalanceRequest)

	return req
}


func (this *GetBalanceRequest) ExchangeMap(values map[string][]string) error {
	var err error
	err = this.exchangeFieldString("playerId", values, &this.PlayerId)
	if err != nil {
		return err
	}

	return nil
}

func (this *GetBalanceRequest) HasErrors() []error {
	valid := validation.Validation{}
	valid.Required(this.PlayerId, "playerId")
	return this.returnErrors(valid)
}
