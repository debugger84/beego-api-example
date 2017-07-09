package errors

import "fmt"

type PlayerUnavailableError struct {
    PlayerId string
}

func (this PlayerUnavailableError) Error() string {
    return fmt.Sprintf("Player with identifier '%s' is unavailable ", this.PlayerId)
}