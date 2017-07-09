package errors

import (
    "fmt"
)

type TournamentNotExistsError struct {
    TournamentId int
}

func (this TournamentNotExistsError) Error() string {
    return fmt.Sprintf("The tournament with identifier '%d' is not exists.", this.TournamentId)
}