package models

type TournamentContribution struct {
	TournamentId  int
	PlayerId      string
	Contribution  int
	BackerContributions []*TournamentContribution
}
