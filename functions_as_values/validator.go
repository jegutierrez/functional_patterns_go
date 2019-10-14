package main

import "log"

// Movement represent an account movement.
type Movement struct {
	ID           int
	Amount       float64
	Fee          float64
	MovementType string
}

type validator func(Movement) bool

// MovementValidator validates the correct form of a movement.
var MovementValidator = map[string]validator{
	"income": func(m Movement) bool {
		if m.Amount <= 0 {
			return false
		}
		if m.Fee <= 0 {
			return false
		}
		return true
	},
	"expense": func(m Movement) bool {
		if m.Amount >= 0 {
			return false
		}
		return true
	},
}

func main() {
	validIncome := Movement{
		ID:           1,
		Amount:       10,
		Fee:          1,
		MovementType: "income",
	}
	validExpense := Movement{
		ID:           2,
		Amount:       -10,
		MovementType: "expense",
	}
	invalidIncomeMov := Movement{
		ID:           3,
		Amount:       10,
		MovementType: "income",
	}

	if !MovementValidator[validIncome.MovementType](validIncome) {
		log.Printf("Invalid movement %d", validIncome.ID)
	}
	if !MovementValidator[validExpense.MovementType](validExpense) {
		log.Printf("Invalid movement %d", validExpense.ID)
	}
	if !MovementValidator[invalidIncomeMov.MovementType](invalidIncomeMov) {
		log.Printf("Invalid movement %d", invalidIncomeMov.ID)
	}
}
