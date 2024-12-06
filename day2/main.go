package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func get_data() ([][]int, error) {
	bytes, err := os.ReadFile("2_data.txt")
	if err != nil {
		return nil, err
	}

	content := string(bytes)
	// 	content := `7 6 4 2 1
	// 1 2 7 8 9
	// 9 7 6 2 1
	// 1 3 2 4 5
	// 8 6 4 4 1
	// 1 3 6 7 9`

	data := [][]int{}
	for _, entry := range strings.Split(content, "\n") {
		nums := strings.Split(entry, " ")

		temp := []int{}
		for _, e := range nums {
			num, err := strconv.ParseInt(e, 10, 64)
			if err != nil {
				return nil, err
			}

			temp = append(temp, int(num))
		}

		data = append(data, temp)
	}

	return data, nil
}

func is_safe(report []int) bool {
	incr := report[0] < report[1]

	for i := 1; i < len(report); i++ {
		step := report[i-1] - report[i]
		stepAbs := math.Abs(float64(step))

		if (step < 0) != incr || stepAbs < 1 || stepAbs > 3 {
			return false
		}
	}
	return true
}

func is_safe_with_dampener(report []int) bool {
	if is_safe(report) {
		return true
	}

	for i := 0; i < len(report); i++ {
		copied := append([]int(nil), report...)

		modifiedReport := append(copied[:i], copied[i+1:]...)

		if is_safe(modifiedReport) {
			return true
		}
	}

	return false
}

func main() {
	data, err := get_data()
	if err != nil {
		panic(err)
	}

	valid_reports := 0

	for _, e := range data {
		if is_safe_with_dampener(e) {
			valid_reports++
		}
	}

	fmt.Println(valid_reports)
}
