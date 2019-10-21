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
	fib := make([]int, n+1, n+2)
	if n < 2 {
		fib = fib[0:2]
	}
	fib[0] = 0
	fib[1] = 1
	for i := 2; i <= n; i++ {
		fib[i] = fib[i-1] + fib[i-2]
	}
	return fib[n]
}
