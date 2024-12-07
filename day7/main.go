package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const TEST_DATA = `190: 10 19
3267: 81 40 27
83: 17 5
156: 15 6
7290: 6 8 6 15
161011: 16 10 13
192: 17 8 14
21037: 9 7 18 13
292: 11 6 16 20`

type Equation struct {
	numbers []int
	result  int
}

func (e *Equation) Equals(operators []string) bool {
	res := e.numbers[0]

	for i := 1; i < len(e.numbers); i++ {
		if operators[i-1] == "*" {
			res *= e.numbers[i]
		} else if operators[i-1] == "+" {
			res += e.numbers[i]
		} else if operators[i-1] == "||" {
			newRes := strconv.Itoa(res) + strconv.Itoa(e.numbers[i])
			intRes, _ := strconv.Atoi(newRes)

			res = intRes
		}
	}

	return e.result == res
}

func permutationsWithReplacement(arr []int, k int) [][]int {
	result := [][]int{}
	generatePermutations(arr, []int{}, k, &result)
	return result
}

func generatePermutations(arr, current []int, k int, result *[][]int) {
	if len(current) == k {
		*result = append(*result, append([]int{}, current...))
		return
	}

	for _, num := range arr {
		generatePermutations(arr, append(current, num), k, result)
	}
}

func operatorsFromIdxs(idxs []int) []string {
	operators := []string{}

	operatorMap := map[int]string{
		0: "+",
		1: "*",
		2: "||",
	}

	for _, idx := range idxs {
		operators = append(operators, operatorMap[idx])
	}

	return operators
}

func (e *Equation) IsPossible() bool {
	if e.Equals(e.genOperators()) {
		return true
	}

	checks := permutationsWithReplacement([]int{0, 1, 2}, len(e.numbers)-1)

	for _, check := range checks {
		if e.Equals(operatorsFromIdxs(check)) {
			return true
		}
	}

	return false
}

func (e *Equation) genOperators() []string {
	operators := []string{}
	for i := 0; i < len(e.numbers)-1; i++ {
		operators = append(operators, "+")
	}

	return operators
}

func parseInput(content string) []Equation {
	equations := []Equation{}

	lines := strings.Split(content, "\n")

	for _, line := range lines {
		splittedLine := strings.Split(line, ": ")

		result, _ := strconv.Atoi(splittedLine[0])

		numbers := []int{}

		for _, num := range strings.Split(splittedLine[1], " ") {
			n, _ := strconv.Atoi(num)
			numbers = append(numbers, n)
		}

		equations = append(equations, Equation{
			result:  result,
			numbers: numbers,
		})
	}

	return equations
}

func main() {
	data, _ := os.ReadFile("data.txt")
	equations := parseInput(string(data))

	sum := 0
	for _, e := range equations {
		if e.IsPossible() {
			sum += e.result
		}
		// fmt.Println(e.IsPossible(), e.result)
	}

	fmt.Println(sum)
}
