package main

import "log"

func main() {
	movements := []AccountMovement{
		{ID: 1, From: "a", To: "b", Amount: 7},
		{ID: 2, From: "c", To: "b", Amount: 14},
		{ID: 3, From: "d", To: "e", Amount: 22},
		{ID: 4, From: "a", To: "f", Amount: 56},
		{ID: 5, From: "g", To: "b", Amount: 8},
		{ID: 6, From: "e", To: "i", Amount: 45},
	}
	var tinyMovements []AccountMovement
	MinFilter(len(movements), 10, func(i int) float64 {
		return movements[i].Amount
	}, func(i int) {
		tinyMovements = append(tinyMovements, movements[i])
	})
	log.Printf("%+v\n", tinyMovements)

	debts := []Debt{
		{ID: 1, Reason: "x", UserID: 4, Amount: 16},
		{ID: 2, Reason: "x", UserID: 2, Amount: 4},
		{ID: 3, Reason: "x", UserID: 1, Amount: 12},
		{ID: 4, Reason: "x", UserID: 3, Amount: 36},
		{ID: 5, Reason: "x", UserID: 1, Amount: 18},
		{ID: 6, Reason: "x", UserID: 3, Amount: 5},
		{ID: 7, Reason: "x", UserID: 4, Amount: 15},
		{ID: 8, Reason: "x", UserID: 2, Amount: 26},
	}
	var tinyDebts []Debt
	MinFilter(len(debts), 10, func(i int) float64 {
		return debts[i].Amount
	}, func(i int) {
		tinyDebts = append(tinyDebts, debts[i])
	})
	log.Printf("%+v\n", tinyDebts)

	var bigMovements []AccountMovement
	Filter(len(movements), func(i int) bool {
		return movements[i].Amount > 20
	}, func(i int) {
		bigMovements = append(bigMovements, movements[i])
	})

	var bigDebts []Debt
	Filter(len(debts), func(i int) bool {
		return debts[i].Amount > 20
	}, func(i int) {
		bigDebts = append(bigDebts, debts[i])
	})

	log.Printf("%+v\n", bigMovements)
}
