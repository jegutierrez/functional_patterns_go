package main

// AccountMovement represents a movement of the user account
type AccountMovement struct {
	ID     int
	From   string
	To     string
	Amount float64
}

// Debt represents a user's debt
type Debt struct {
	ID     int
	UserID int
	Reason string
	Amount float64
}

// MinFilter generic function to filter a slice on min Amount.
// Params:
// l 			-> slice length
// min 			-> min amount
// amountGetter -> closure to get the Amount value
// appender 	-> closure to append to the original slice
func MinFilter(l int, min float64, amountGetter func(int) float64, appender func(int)) {
	for i := 0; i < l; i++ {
		if amountGetter(i) < min {
			appender(i)
		}
	}
}

// Filter generic function to filter a slice on a given predicate.
// Params:
// l 			-> slice length
// amountGetter -> closure to get the Amount value
// appender 	-> closure to append to the original slice
func Filter(l int, predicate func(int) bool, appender func(int)) {
	for i := 0; i < l; i++ {
		if predicate(i) {
			appender(i)
		}
	}
}
