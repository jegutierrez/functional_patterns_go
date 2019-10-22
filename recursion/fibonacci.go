package main

// FibonacciRecursive implementation using recursion.
func FibonacciRecursive(n int) int {
	if n <= 1 {
		return n
	}
	return FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
}

// FibonacciFor implementation using for loops.
func FibonacciFor(n int) int {
	if n <= 1 {
		return n
	}

	var n2, n1 int = 0, 1

	for i := 2; i < n; i++ {
		n2, n1 = n1, n1+n2
	}

	return n2 + n1
}
