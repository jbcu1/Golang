package main

import (
	"fmt"
	"net"
)

func HighestRank(nums []int) int {
	num := 0
	freq := 0
	m := make(map[int]int)
	for _, l := range nums {
		count := 0
		for _, l1 := range nums {
			if l == l1 {
				count++
			}
		}
		m[l] = count
	}

	for i, n := range m {
		if n > num {
			num = n
			freq = i
		}

	}
	return freq

}

func main() {

	fmt.Println(HighestRank([]int{172, 38, 205, 49, 94, 251, 89, 49, 105, 236, 174, 140, 198, 86, 160}))
}

func Is_valid_ip(ip string) bool {
	res := net.ParseIP(ip)
	if res == nil {
		return false
	}
	return true
}
