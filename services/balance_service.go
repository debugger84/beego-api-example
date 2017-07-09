package services

import (
	"github.com/astaxie/beego/logs"
	serviceErrors "tournamentAPI/services/errors"
	"tournamentAPI/storage/interfaces"
	"tournamentAPI/storage"
	"tournamentAPI/requests"
	"tournamentAPI/models"
)

type BalanceService struct {
	log *logs.BeeLogger
	playerStorage interfaces.PlayerStorageInterface
}


func NewBalanceService() *BalanceService {
	var service = new(BalanceService)

	log := logs.NewLogger(1)
	log.SetLogger("file", `{"filename":"logs/balance.log"}`)
	service.log = log
	service.playerStorage = storage.NewPlayerRedisStorage()

	return service
}

//Getting of player balance
func (this *BalanceService) GetPlayerBalance(playerId string) (int, error) {
	player := this.playerStorage.GetPlayer(playerId)
	if player == nil {
		return 0, &serviceErrors.PlayerUnavailableError{PlayerId: playerId}
	}

	return player.Balance, nil
}


func (this *BalanceService) IncreasePlayerBalance(request *requests.ChangeBalanceRequest) (error) {
	var player *models.Player
	playerId := request.PlayerId
	player = this.playerStorage.GetPlayer(playerId)
	if player == nil {
		player = &models.Player{
			PlayerId: request.PlayerId,
			Balance: 0,
		}
	}

	player.Balance += request.Points

	err := this.playerStorage.SavePlayer(player)
	if err != nil {
		this.log.Error("The error has been occurred during saving the player's balance: %s", err.Error())
	}

	return err
}

func (this *BalanceService) DecreasePlayerBalance(request *requests.ChangeBalanceRequest) (error) {
	var player *models.Player
	playerId := request.PlayerId
	player = this.playerStorage.GetPlayer(playerId)
	if player == nil {
		return &serviceErrors.PlayerUnavailableError{PlayerId: playerId}
	}

	player.Balance -= request.Points

	if player.Balance < 0 {
		this.log.Error("Trying to change balance to negative value for the player: %s", playerId)
		return &serviceErrors.BalanceIsNegativeError{PlayerId: playerId, Balance:player.Balance}
	}
	err := this.playerStorage.SavePlayer(player)
	if err != nil {
		this.log.Error("The error has been occurred during saving the player's balance: %s", err.Error())
	}

	return err
}