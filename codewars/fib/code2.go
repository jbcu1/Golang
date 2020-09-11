package main

import "math"

func fib() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

func main() {

	f := fib()
	println(f(), f(), f(), f(), f())
}

func LastFibDigit(n int) int {
	i := math.Pow(5.0, 0.5)
	left := (1 + i) / 2
	right := (1 - i) / 2
	f := (int((math.Pow(left, float64(n))-math.Pow(right, float64(n)))/i) % 10)
	return f
}
