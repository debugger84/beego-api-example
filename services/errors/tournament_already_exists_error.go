package errors

import (
    "fmt"
)

type TournamentAlreadyExistsError struct {
    TournamentId int
}

func (this TournamentAlreadyExistsError) Error() string {
    return fmt.Sprintf("The tournament with identifier '%d' is already exists.", this.TournamentId)
}