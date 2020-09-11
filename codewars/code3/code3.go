
import (
	"fmt"
	"math"
)

func MxDifLg(a1 []string, a2 []string) int {
	c := 0
	var max []int
	if len(a1) == 0 || len(a2) == 0 {
		return -1
	}

	for k := range a1 {
		for j := range a2 {
			c := math.Abs(len(a1[k]) - len(a2[j]))
			max = append(max, c)

		}
	}
	for _, o := range max {

		if o > c {
			c = o
		}
	}

	return c
}


func GetCount(str string) (count int) {
	for _, c := range str {
		switch c {
		case 'a', 'e', 'i', 'o', 'u':
			count++
		}
	}
	return count
}
