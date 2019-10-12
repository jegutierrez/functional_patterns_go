package main

// Parentheses Balancing
// Write a function which verifies the balancing of parentheses
// in a string. For example, the function should return true for
// the following strings:

//  - (if (zero? x) max (/ 1 x))
//  - I told him (that it’s not (yet) done). (But he wasn’t listening)

// The function should return false for the following strings:
//  - :-)
//  - ())(

// BalanceRecursive tail recursive implementation.
func BalanceRecursive(s string) bool {
	return balanceHelper(s, 0, 0)
}

func balanceHelper(s string, pointer, count int) bool {
	if len(s) == pointer {
		return count == 0
	}
	if s[pointer] == '(' {
		return balanceHelper(s, pointer+1, count+1)
	} else if s[pointer] == ')' {
		if count == 0 {
			return false
		}
		return balanceHelper(s, pointer+1, count-1)
	}
	return balanceHelper(s, pointer+1, count)
}

// BalanceGoStyle iterative Go style implementation.
func BalanceGoStyle(s string) bool {
	var count int
	for _, c := range s {
		if c == '(' {
			count++
		} else if c == ')' {
			if count == 0 {
				return false
			}
			count--
		}
	}
	return count == 0
}

func main() {}
