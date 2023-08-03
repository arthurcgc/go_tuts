package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

func dfs(root *tree.Tree, ch chan int) {
	if root == nil {
		return
	}

	dfs(root.Left, ch)
	ch <- root.Value
	dfs(root.Right, ch)
}

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	dfs(t, ch)
	close(ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)
	for {
		v1, ok := <-ch1
		if !ok {
			break
		}
		v2, ok := <-ch2
		if !ok {
			break
		}
		fmt.Printf("%d - %d\n", v1, v2)
		if v1 != v2 {
			return false
		}
	}
	return true
}

func main() {
	t1 := tree.New(1)
	t2 := tree.New(1)
	t3 := tree.New(2)

	fmt.Println(Same(t1, t2))
	fmt.Println(Same(t1, t3))

}
