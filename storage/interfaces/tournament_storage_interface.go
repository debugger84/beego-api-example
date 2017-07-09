package interfaces

import "tournamentAPI/models"

type TournamentStorageInterface interface {
    GetTournament(tournamentId int) *models.Tournament
    SaveTournament(tournament *models.Tournament) error
    RemoveTournament(tournament *models.Tournament) error
}
