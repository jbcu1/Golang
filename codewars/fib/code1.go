package main

import "math"

func FindOdd(seq []int) int {
	a := 0
	for _, num := range seq {
		count := 0
		for _, num1 := range seq {
			if num == num1 {
				count++
			}
		}
		if math.Mod(float64(count), 2.0) != 0 {
			a = num
		}
	}
	return a
}

func FindOdd1(seq []int) int {
	res := 0
	for _, x := range seq {
		res ^= x
	}
	return res
}
