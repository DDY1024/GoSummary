package main

import (
	"sort"
)

// 二分查找的方式
// 我们完全可以 O(n) 复杂度进行预处理
// 左右各自进行相应的递推
func platesBetweenCandlesOne(s string, queries [][]int) []int {
	n, m := len(s), len(queries)
	ans := make([]int, m)
	sum := make([]int, n)
	pos := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if s[i] == '|' {
			sum[i]++
			pos = append(pos, i)
		}
		if i-1 >= 0 {
			sum[i] += sum[i-1]
		}
	}

	var calSum = func(i, j int) int {
		ret := sum[j]
		if i-1 >= 0 {
			ret -= sum[i-1]
		}
		return ret
	}

	// 二分查找
	for i := 0; i < m; i++ {
		l, r := queries[i][0], queries[i][1]
		idx1 := sort.Search(len(pos), func(i int) bool {
			return pos[i] >= l
		})
		if idx1 >= len(pos) || pos[idx1] > r {
			continue
		}
		idx2 := sort.Search(len(pos), func(i int) bool {
			return pos[i] > r
		})
		idx2--
		ans[i] = pos[idx2] - pos[idx1] + 1 - calSum(pos[idx1], pos[idx2])
	}
	return ans
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// O(n) 算法复杂度进行相应的递推
func platesBetweenCandles(s string, queries [][]int) []int {
	n, m := len(s), len(queries)
	sum := make([]int, n)
	dl, dr := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		dr[i] = -1
		if s[i] == '|' {
			sum[i]++
			dr[i] = i
		}
		if i-1 >= 0 {
			sum[i] += sum[i-1]
			dr[i] = maxInt(dr[i], dr[i-1]) // max 操作
		}
	}
	for i := n - 1; i >= 0; i-- {
		dl[i] = 0x3f3f3f3f
		if s[i] == '|' {
			dl[i] = i
		}
		if i+1 < n {
			dl[i] = minInt(dl[i], dl[i+1]) // min 操作
		}
	}

	var calSum = func(i, j int) int {
		ret := sum[j]
		if i-1 >= 0 {
			ret -= sum[i-1]
		}
		return ret
	}

	ans := make([]int, m)
	for i := 0; i < m; i++ {
		l, r := queries[i][0], queries[i][1]
		ll, rr := dl[l], dr[r]
		if ll > rr {
			continue
		}
		ans[i] = rr - ll + 1 - calSum(ll, rr)
	}
	return ans
}

func main() {
	// var mux sync.Mutex
}
