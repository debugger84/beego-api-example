package models

type Tournament struct {
	TournamentId  int
	Deposit       int
	Contributions []*TournamentContribution
}
