package main

import (
	"fmt"
)

const mod = int64(1e9 + 7)

func main() {

	n, cnt := 0, 0
	fmt.Scan(&n, &cnt)
	arr := make([]int, 0, n)
	for i := 0; i < cnt; i++ {
		x, p := 0, 0
		fmt.Scan(&x, &p)
		for j := 0; j < p; j++ {
			arr = append(arr, x)
		}
	}
	fmt.Scan(&cnt)
	for i := 0; i < cnt; i++ {
		start, end := 0, 0
		fmt.Scan(&start, &end)
		sub := arr[start-1 : end]
		ans := int64(0)
		for i := 0; i < len(sub); i++ {
			if i != 0 && sub[i] == sub[i-1] {
				continue
			}
			if sub[i] == 1 {
				ans++
				continue
			}
			ans += int64(len(sub) - i)

		}
		fmt.Println(ans % mod)
	}

}
