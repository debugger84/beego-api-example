package services

import (
	"github.com/astaxie/beego/logs"
	serviceErrors "tournamentAPI/services/errors"
	"tournamentAPI/storage/interfaces"
	"tournamentAPI/storage"
	"tournamentAPI/requests"
	"tournamentAPI/models"
    "errors"
)

type TournamentService struct {
	log               *logs.BeeLogger
	tournamentStorage interfaces.TournamentStorageInterface
	balanceService    *BalanceService
}


func NewTournamentService() *TournamentService {
	var service = new(TournamentService)

	log := logs.NewLogger(1)
	log.SetLogger("file", `{"filename":"logs/tournament.log"}`)
	service.log = log
	service.tournamentStorage = storage.NewTournamentRedisStorage()
	service.balanceService = NewBalanceService()

	return service
}

//Creates tournament
func (this *TournamentService) AddTournament(request *requests.AnnounceTournamentRequest) (error) {
	tournament := this.tournamentStorage.GetTournament(request.TournamentId)
	if tournament != nil {
		return &serviceErrors.TournamentAlreadyExistsError{TournamentId: request.TournamentId}
	}

	tournament = &models.Tournament{
		TournamentId: request.TournamentId,
		Deposit: request.Deposit,
	}

	this.tournamentStorage.SaveTournament(tournament)

	return  nil
}

//Joins players to tournament
func (this *TournamentService) JoinTournament(request *requests.JoinTournamentRequest) (error) {

	tournament := this.tournamentStorage.GetTournament(request.TournamentId)
	if tournament == nil {
		return &serviceErrors.TournamentNotExistsError{TournamentId: request.TournamentId}
	}

    if this.playerInTournament(tournament, request.PlayerId) {
        return errors.New("The player has been joined already ")
    }

    playerContribution, backerContribution := this.calcPlayerAndBackerParts(
        tournament.Deposit,
        len(request.BackerIds),
    )

    contribution := &models.TournamentContribution{
        PlayerId: request.PlayerId,
        Contribution: playerContribution,
        TournamentId: request.TournamentId,
    }

    changeBalReq := &requests.ChangeBalanceRequest{
        PlayerId: request.PlayerId,
        Points: playerContribution,
    }
    err := this.balanceService.DecreasePlayerBalance(changeBalReq)
    if err != nil {
        return err
    }

    backers, joinErr := this.joinBackers(tournament, request.BackerIds, backerContribution)
    if joinErr != nil {
        this.balanceService.IncreasePlayerBalance(changeBalReq)
        return joinErr
    }

    contribution.BackerContributions = backers

    tournament.Contributions = append(tournament.Contributions, contribution)

	return  this.tournamentStorage.SaveTournament(tournament)
}

//Closes tournament and gives prizes to winners
func (this *TournamentService) CloseTournament(request *requests.CloseTournamentRequest) (error) {

    tournament := this.tournamentStorage.GetTournament(request.TournamentId)
    if tournament == nil {
        return &serviceErrors.TournamentNotExistsError{TournamentId: request.TournamentId}
    }

    err := this.givePrizes(tournament, request.Winners)
    if err != nil {
        this.log.Error("Can't set prizes to winners")
        return errors.New("Can't set prizes to winners")
    }

    err = this.tournamentStorage.RemoveTournament(tournament)
    if err != nil {
        this.log.Error("Can't remove tournament")
        return errors.New("Can't remove tournament")
    }

    return nil
}

func (this *TournamentService) givePrizes(tournament *models.Tournament, winners []*requests.Winner) (error) {
    var playerPrize int
    var backerPrize int
    var changeBalReq *requests.ChangeBalanceRequest
    for _, winner := range winners {
        for _, contribution := range tournament.Contributions {
            if winner.PlayerId == contribution.PlayerId {
                playerPrize, backerPrize = this.calcPlayerAndBackerParts(
                    winner.Prize,
                    len(contribution.BackerContributions),
                )
                changeBalReq = &requests.ChangeBalanceRequest{
                    PlayerId: winner.PlayerId,
                    Points: playerPrize,
                }
                this.balanceService.IncreasePlayerBalance(changeBalReq)
                for _, backer := range contribution.BackerContributions {
                    changeBalReq = &requests.ChangeBalanceRequest{
                        PlayerId: backer.PlayerId,
                        Points: backerPrize,
                    }
                    this.balanceService.IncreasePlayerBalance(changeBalReq)
                }
            }
        }
    }

    return nil
}

func (this *TournamentService) calcPlayerAndBackerParts(wholeSum int, backerCount int) (int, int) {
    playersCount := backerCount + 1
    if playersCount == 1 {
        return wholeSum, 0
    }
    contribution := int(float64(wholeSum / playersCount) + float64(0.5))

    playerContribution := wholeSum - (contribution * (playersCount -1))

    return playerContribution, contribution
}

func (this *TournamentService) playerInTournament(tournament *models.Tournament, playerId string) bool {
    for _, v := range tournament.Contributions {
        if v.PlayerId == playerId {
            return true
        }
    }

    return false
}

func (this *TournamentService) joinBackers(tournament *models.Tournament, backersIds []string, contribution int) ([]*models.TournamentContribution, error) {
    var backers []*models.TournamentContribution
    var changeBalReq *requests.ChangeBalanceRequest
    var err error

    for _, v := range backersIds {
        changeBalReq = &requests.ChangeBalanceRequest{
            PlayerId: v,
            Points: contribution,
        }
        err = this.balanceService.DecreasePlayerBalance(changeBalReq)
        if err != nil {
            this.rollbackBackers(backers)
            return nil, err
        }
        backers = append(backers, &models.TournamentContribution{
            PlayerId: v,
            Contribution: contribution,
            TournamentId: tournament.TournamentId,
        })
    }

    return backers, nil
}

func (this *TournamentService) rollbackBackers(backers []*models.TournamentContribution) {
    var changeBalReq *requests.ChangeBalanceRequest
    var err error

    for _, v := range backers {
        changeBalReq = &requests.ChangeBalanceRequest{
            PlayerId: v.PlayerId,
            Points:   v.Contribution,
        }
        err = this.balanceService.IncreasePlayerBalance(changeBalReq)
        if err != nil {
            this.log.Error("Rollback of balance charging for %d points is broken for PlayerId %s", v.PlayerId, v.Contribution)
        }
    }
}
