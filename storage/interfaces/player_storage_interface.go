package interfaces

import "tournamentAPI/models"

//@see PlayerRedisStorage
type PlayerStorageInterface interface {
    GetPlayer(playerId string) *models.Player
    SavePlayer(player *models.Player) error
}
