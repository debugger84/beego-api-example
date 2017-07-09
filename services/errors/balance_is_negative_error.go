package errors

import "fmt"

type BalanceIsNegativeError struct {
    PlayerId string
    Balance int
}

func (this BalanceIsNegativeError) Error() string {
    return fmt.Sprintf("Balance can't be negative for the player with identifier '%s'. Current balance is %d", this.PlayerId, this.Balance)
}