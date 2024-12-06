package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func get_data() ([][2]int, error) {
	bytes, err := os.ReadFile("1_1data.txt")
	if err != nil {
		return nil, err
	}

	content := string(bytes)

	data := [][2]int{}
	for _, entry := range strings.Split(content, "\n") {
		nums := strings.Split(entry, "   ")

		num1, err1 := strconv.ParseInt(nums[0], 10, 64)
		if err1 != nil {
			return nil, err1
		}

		num2, err2 := strconv.ParseInt(nums[1], 10, 64)
		if err2 != nil {
			return nil, err2
		}

		data = append(data, [2]int{int(num1), int(num2)})
	}

	return data, nil
}

func main() {
	data, err := get_data()
	if err != nil {
		panic(err)
	}

	left := []int{}
	for _, e := range data {
		left = append(left, e[0])
	}

	right := []int{}
	for _, e := range data {
		right = append(right, e[1])
	}

	// sort.Ints(left)
	// sort.Ints(right)
	counts := map[int]int{}

	for _, e := range right {
		val, ok := counts[e]
		if ok {
			counts[e] = val + 1
		} else {
			counts[e] = 1
		}
	}

	sum := 0

	for _, e := range left {
		val, ok := counts[e]

		if ok {
			sum = sum + e*val
		}
	}

	fmt.Println(sum)
}
