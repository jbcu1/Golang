package main

import (
	"fmt"
)

func main() {
	var n uint64
	n = 7
	var sumfuc uint64
	for i:=1;i<=int(n);i++{
		sumfuc+=factorial(uint64(i))
	}
	fmt.Println(sumfuc, factorial(n))


	l:=(1/float64(factorial(n)))*float64(sumfuc)
	k,_:=fmt.Printf("%.4f",l)
	fmt.Println(float64(k))
}

func factorial(n uint64) (n1 uint64) {
	if n > 0 {
		n1 = n * factorial(n-1)

		return n1
	}
	return 1
}
