package main

import (
	"fmt"
	"golang.org/x/tour/tree"
)

type Tree struct {
	Left  *Tree
	Value int
	Right *Tree
}

func Walk(t *tree.Tree, ch chan int) {
	_walk(t, ch)
	close(ch)
}

func _walk(t *tree.Tree, ch chan int) {
	if t != nil {
		_walk(t.Left, ch)
		ch <- t.Value
		_walk(t.Right, ch)
	}
}

func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)
	for i := range ch1 {
		if i != <-ch2 {
			return false
		}
	}
	return true
}

func main() {
	ch := make(chan int)
	ch1 := make(chan int)
	go Walk(tree.New(1), ch)
	go Walk(tree.New(2), ch1)
	for v := range ch {
		fmt.Println(v)
	}
	fmt.Println()
	for x := range ch1 {
		fmt.Println(x)
	}
	fmt.Println()
	fmt.Println(Same(tree.New(1), tree.New(2)))
	fmt.Println()
	fmt.Println(Same(tree.New(1), tree.New(1)))
}
