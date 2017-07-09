package storage

import (
    "github.com/go-redis/redis"
    "tournamentAPI/models"
    "tournamentAPI/basic"
    "github.com/astaxie/beego/logs"
    "strconv"
)

type PlayerRedisStorage struct {
    client *redis.Client
    log *logs.BeeLogger
}


func NewPlayerRedisStorage() *PlayerRedisStorage {
    var storage = new(PlayerRedisStorage)

    client := basic.GetRedisClient()

    storage.client = client
    log := logs.NewLogger(1)
    log.SetLogger("file", `{"filename":"logs/redis.log"}`)
    storage.log = log

    return storage
}

func (this *PlayerRedisStorage) GetPlayer(playerId string) *models.Player {
    balanceRes, err := this.client.Get(this.getPlayerBalanceKey(playerId)).Result()
    if err != nil {
        this.log.Info("Trying to access unavailable player " + playerId)
        return nil
    }

    balance, convErr := strconv.Atoi(balanceRes)
    if convErr != nil {
        this.log.Error("Converting of user balance is wrong for player " + playerId)
        return nil
    }

    return &models.Player{
        PlayerId: playerId,
        Balance: balance,
    }
}

func (this *PlayerRedisStorage) SavePlayer(player *models.Player) error {
    key := this.getPlayerBalanceKey(player.PlayerId)
    err := this.client.Set(key, player.Balance, 0).Err()

    return err
}

func (this *PlayerRedisStorage) getPlayerBalanceKey(playerId string) string {
    return "player_balance:" + playerId
}