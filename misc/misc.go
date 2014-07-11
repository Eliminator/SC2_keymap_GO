package misc

import (
	"fmt"
	"time"
)

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func UniqInt(col []int) []int {
	m := map[int]struct{}{}
	for _, v := range col {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
		}
	}
	list := make([]int, len(m))

	i := 0
	for v := range m {
		list[i] = v
		i++
	}
	return list
}

func IntMin(a int, b int) int {
	if a > b {
		return b
	} else {
		return a
	}

}
func IntMax(a int, b int) int {
	if a < b {
		return b
	} else {
		return a
	}

}

func IntInArray(a int, arr []int) bool {
	for _, b := range arr {
		if b == a {
			return true
		}
	}
	return false
}

func ByteInArray(a byte, arr []byte) bool {
	for _, b := range arr {
		if b == a {
			return true
		}
	}
	return false
}
