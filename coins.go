package main

// Counting Change
// Write a recursive function that counts how many different ways
// you can make change for an amount, given a list of coin denominations.
// For example, there are 3 ways to give change for 4 if you have coins with
// denomination 1 and 2: 1+1+1+1, 1+1+2, 2+2.

// CoinsChangeRecursive recursive implementation.
func CoinsChangeRecursive(amount int, coins []int) int {
	if amount == 0 {
		return 1
	} else if amount > 0 && len(coins) > 0 {
		return CoinsChangeRecursive(amount-coins[0], coins) + CoinsChangeRecursive(amount, coins[1:])
	} else {
		return 0
	}
}

// CoinsChangeGoStyle iterative implementation.
func CoinsChangeGoStyle(amount int, coins []int) int {
	var table = make([]int, amount+1, amount+1)
	table[0] = 1

	for i := 0; i < len(coins); i++ {
		for j := coins[i]; j <= amount; j++ {
			table[j] += table[j-coins[i]]
		}
	}
	return table[amount]
}
