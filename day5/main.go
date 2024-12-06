package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

const content = `47|53
97|13
97|61
97|47
75|29
61|13
75|53
29|13
97|29
53|29
61|53
97|53
61|29
47|13
75|47
97|75
47|61
75|61
47|29
75|13
53|13

75,47,61,53,29
97,61,53,29,13
75,29,13
75,97,47,61,53
61,13,29
97,13,75,29,47`

func parseInput(content string) (map[string][]string, [][]string) {
	splitted := strings.Split(content, "\n\n")

	rulesStr := splitted[0]
	rules := map[string][]string{}

	for _, line := range strings.Split(rulesStr, "\n") {
		splittedLine := strings.Split(line, "|")

		before, after := splittedLine[0], splittedLine[1]

		val, ok := rules[before]
		if ok {
			rules[before] = append(val, after)
		} else {
			rules[before] = []string{after}
		}
	}

	updatesStr := splitted[1]
	updates := [][]string{}

	for _, line := range strings.Split(updatesStr, "\n") {
		splittedLine := strings.Split(line, ",")
		updates = append(updates, splittedLine)
	}

	return rules, updates
}

func getMiddleEl[T any](arr []T) T {
	middle := len(arr) / 2
	return arr[middle]
}

func isValidUpdate(update []string, rules map[string][]string) bool {
	for i := 0; i < len(update); i++ {
		currentPage := update[i]

		for j := i + 1; j < len(update); j++ {
			checkPage := update[j]

			seq := rules[currentPage]

			if !slices.Contains(seq, checkPage) {
				return false
			}
		}

	}

	return true
}

func reorderUpdate(update []string, rules map[string][]string) []string {
	resUpdate := []string{}

	actualRules := map[string][]string{}

	for key, value := range rules {
		if slices.Contains(update, key) {
			actualRules[key] = value
		}
	}

	for len(actualRules) > 0 {
	outer:
		for key := range actualRules {

		inner:
			for innerKey, innerValue := range actualRules {
				if key == innerKey {
					continue inner
				}

				if slices.Contains(innerValue, key) {
					continue outer
				}
			}

			fmt.Println(actualRules)
			resUpdate = append(resUpdate, key)
			delete(actualRules, key)
		}
	}

	return resUpdate
}

func main() {
	data, _ := os.ReadFile("5_data.txt")
	rules, updates := parseInput(string(data))

	sum := 0

	for _, update := range updates {
		if !isValidUpdate(update, rules) {
			reorderedUpdate := reorderUpdate(update, rules)
			middle := getMiddleEl(reorderedUpdate)

			val, _ := strconv.Atoi(middle)
			sum += val
		}
	}

	fmt.Println(sum)
}
