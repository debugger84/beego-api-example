package storage

import (
    "github.com/go-redis/redis"
    "tournamentAPI/models"
    "tournamentAPI/basic"
    "github.com/astaxie/beego/logs"
    "strconv"
)

type TournamentRedisStorage struct {
    client *redis.Client
    log    *logs.BeeLogger
}

func NewTournamentRedisStorage() *TournamentRedisStorage {
    var storage = new(TournamentRedisStorage)

    client := basic.GetRedisClient()

    storage.client = client
    log := logs.NewLogger(1)
    log.SetLogger("file", `{"filename":"logs/redis.log"}`)
    storage.log = log

    return storage
}

func (this *TournamentRedisStorage) GetTournament(tournamentId int) *models.Tournament {
    var contributions []*models.TournamentContribution
    var backerContributions []*models.TournamentContribution
    depositRes, err := this.client.Get(this.getTournamentDepositKey(tournamentId)).Result()
    if err != nil {
        this.log.Info("Trying to access unavailable tournament " + strconv.Itoa(tournamentId))
        return nil
    }

    deposit, convErr := strconv.Atoi(depositRes)
    if convErr != nil {
        this.log.Error("Converting of user deposit is wrong for tournament " + strconv.Itoa(tournamentId))
        return nil
    }

    contributionsRes, err := this.client.HGetAll(this.getPlayerContributionKey(tournamentId)).Result()
    if err == nil {
        contributions = this.convertMapToContributions(contributionsRes, tournamentId)
        for _, contribution := range contributions {
            contributionsRes, err = this.client.HGetAll(this.getBackerContributionKey(tournamentId, contribution.PlayerId)).Result()
            if err == nil {
                backerContributions = this.convertMapToContributions(contributionsRes, tournamentId)
                contribution.BackerContributions = backerContributions
            }
        }
    }

    return &models.Tournament{
        TournamentId:  tournamentId,
        Deposit:       deposit,
        Contributions: contributions,
    }
}

//@TODO Need to emulate transactions
func (this *TournamentRedisStorage) SaveTournament(tournament *models.Tournament) error {
    err := this.saveTournamentData(tournament)
    if err != nil {
        return err
    }

    err = this.saveTournamentPlayersContributions(tournament)
    if err != nil {
        return err
    }

    err = this.saveTournamentBackersContributions(tournament)
    if err != nil {
        return err
    }

    return nil
}

func (this *TournamentRedisStorage) RemoveTournament(tournament *models.Tournament) error {
    tournamentId := tournament.TournamentId
    key := this.getTournamentDepositKey(tournamentId)
    this.client.Del(key)

    key = this.getPlayerContributionKey(tournamentId)
    this.client.Del(key)

    for _, contribution := range tournament.Contributions {
        key = this.getBackerContributionKey(tournamentId, contribution.PlayerId)
        this.client.Del(key)
    }

    return nil
}

func (this *TournamentRedisStorage) saveTournamentData(tournament *models.Tournament) error {
    tournamentId := tournament.TournamentId
    key := this.getTournamentDepositKey(tournamentId)
    err := this.client.Set(key, tournament.Deposit, 0).Err()
    if err != nil {
        this.log.Error("Tournament deposit can't be saved " + strconv.Itoa(tournamentId))
        return err
    }

    return nil
}

func (this *TournamentRedisStorage) saveTournamentPlayersContributions(tournament *models.Tournament) error {
    tournamentId := tournament.TournamentId

    if len(tournament.Contributions) == 0 {
        return nil
    }
    key := this.getPlayerContributionKey(tournamentId)
    contributionsMap := this.convertContributionsToMap(tournament.Contributions)
    err := this.client.HMSet(key, contributionsMap).Err()
    if err != nil {
        this.log.Error("Tournament contributions can't be saved " + strconv.Itoa(tournamentId))
        return err
    }

    return nil
}

func (this *TournamentRedisStorage) saveTournamentBackersContributions(tournament *models.Tournament) error {
    tournamentId := tournament.TournamentId

    for _, playerContrib := range tournament.Contributions {
        if len(playerContrib.BackerContributions) == 0 {
            continue
        }
        key := this.getBackerContributionKey(tournamentId, playerContrib.PlayerId)
        contributionsMap := this.convertContributionsToMap(playerContrib.BackerContributions)
        err := this.client.HMSet(key, contributionsMap).Err()
        if err != nil {
            this.log.Error("Tournament contributions of backer can't be saved " + strconv.Itoa(tournamentId))
            return err
        }
    }
    return nil
}

func (this *TournamentRedisStorage) convertMapToContributions(contributionsMap map[string]string, tournamentId int) []*models.TournamentContribution {
    var contributions []*models.TournamentContribution
    var contribution int
    var err error
    //contributions = make([]*models.TournamentContribution, len(contributionsMap))
    for playerId, contributionRes := range contributionsMap {
        contribution, err = strconv.Atoi(contributionRes)
        if err != nil {
            this.log.Error("Contribution value cant't be converted to int")
        }
        contributions = append(contributions, &models.TournamentContribution{
            TournamentId: tournamentId,
            Contribution: contribution,
            PlayerId:     playerId,
        })
    }

    return contributions
}

func (this *TournamentRedisStorage) convertContributionsToMap(contributions []*models.TournamentContribution) map[string]interface{} {
    var result map[string]interface{}
    result = make(map[string]interface{}, len(contributions))
    for _, v := range contributions {
        result[v.PlayerId] = v.Contribution
    }

    return result
}

func (this *TournamentRedisStorage) getTournamentDepositKey(tournamentId int) string {
    return "tournament_deposit:" + strconv.Itoa(tournamentId)
}

func (this *TournamentRedisStorage) getPlayerContributionKey(tournamentId int) string {
    return "tournament_player_contribution:" + strconv.Itoa(tournamentId)
}

func (this *TournamentRedisStorage) getBackerContributionKey(tournamentId int, playerId string) string {
    return "tournament_backer_contribution:" + strconv.Itoa(tournamentId) + ":" + playerId
}
