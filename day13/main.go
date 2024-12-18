package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
)

const test_data = `Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=8400, Y=5400

Button A: X+26, Y+66
Button B: X+67, Y+21
Prize: X=12748, Y=12176

Button A: X+17, Y+86
Button B: X+84, Y+37
Prize: X=7870, Y=6450

Button A: X+69, Y+23
Button B: X+27, Y+71
Prize: X=18641, Y=10279`

type Formula struct {
	buttonA [2]int
	buttonB [2]int

	prize [2]int
}

func (f *Formula) FindMinToken() int {
	tokens := []int{}

	for a := 0; a < 100; a++ {
		for b := 0; b < 100; b++ {
			xRes := f.buttonA[0]*a + f.buttonB[0]*b
			yRes := f.buttonA[1]*a + f.buttonB[1]*b

			if xRes == f.prize[0] && yRes == f.prize[1] {
				tokens = append(tokens, a*3+b)
			}
		}
	}

	if len(tokens) < 1 {
		return 0
	}

	minToken := int(math.Inf(1))
	for _, t := range tokens {
		if minToken > t {
			minToken = t
		}
	}

	return minToken
}

func (f *Formula) BumpPrize() {
	f.prize[0] += 10000000000000
	f.prize[1] += 10000000000000
}

func extractNumbers(line string) [2]int {
	nums := []int{}

	splitted := strings.Split(line, ": ")
	numbers := strings.Split(splitted[1], ", ")

	for _, num := range numbers {
		splittedNum := strings.Split(num, "+")
		if len(splittedNum) < 2 {
			splittedNum = strings.Split(num, "=")
		}

		n, _ := strconv.Atoi(splittedNum[1])
		nums = append(nums, n)
	}

	return [2]int(nums)
}

func NewFormula(block string) Formula {
	lines := strings.Split(block, "\n")

	buttons := [][2]int{}
	for _, button := range lines[:2] {
		buttons = append(buttons, extractNumbers(button))
	}

	prize := extractNumbers(lines[2])

	return Formula{
		buttonA: buttons[0],
		buttonB: buttons[1],

		prize: prize,
	}
}

func main() {
	data, _ := os.ReadFile("data.txt")
	blocks := strings.Split(string(data), "\n\n")

	sum := 0

	done := 0

	var wg sync.WaitGroup

	maxGoroutines := 16
	guard := make(chan struct{}, maxGoroutines)

	for _, block := range blocks {
		wg.Add(1)

		guard <- struct{}{}

		go func(b string) {
			defer wg.Done()

			formula := NewFormula(b)
			// formula.BumpPrize()
			sum += formula.FindMinToken()

			done++
			fmt.Printf("%v/%v %v\n", done, len(blocks), sum)

			<-guard
		}(block)
	}

	wg.Wait()
	fmt.Println(sum)
}
