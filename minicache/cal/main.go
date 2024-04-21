package main

type ListNode struct {
	Val  int
	Next *ListNode
}
type Node struct {
	Val    int
	Next   *Node
	Random *Node
}
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func combinationSum3(k int, n int) [][]int {
	ans := [][]int{}
	path := []int{}
	var dfs func(sum, i, cnt int)
	dfs = func(sum, x, cnt int) {
		if sum > n {
			return
		}
		if sum == n && cnt != k {
			return
		}
		if sum == n && cnt == k {
			ans = append(ans, append([]int{}, path...))
			return
		}
		for i := x; i <= 9; i++ {
			path = append(path, i)
			dfs(sum+i, i+1, cnt+1)
			path = path[:len(path)-1]
		}
	}
	dfs(0, 1, 0)
	return ans
}

func main() {
	//a := []int{3, 2, 1, 5, 6, 4}
	//b := []int{1}
	//s := "000110"
	//grid := [][]int{{1, 0}, {1, 2}, {1, 3}}
	//grid := [][]int{{6, 10}, {5, 15}}
	//fmt.Println(minimumAddedCoins(a, 19))
	//a, b := 1, 5
	//fmt.Println(isValidSerialization(s))

	//fmt.Println(findKthLargest(a, 2))
}
func min(x int, y int) int {
	if x <= y {
		return x
	} else {
		return y
	}
}

func max(x int, y int) int {
	if x <= y {
		return y
	} else {
		return x
	}
}
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
