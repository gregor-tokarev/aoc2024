package main

import (
	"fmt"
	"os"
	"strings"
)

func isXmas(c rune) bool {
	return c == 'X' || c == 'M' || c == 'A' || c == 'S'
}

type Letter struct {
	L   rune
	Pos [2]int
}

type LetterMap struct {
	letters [][]Letter
}

func (lm *LetterMap) Len() int {
	return len(lm.letters)
}

func (lm *LetterMap) Get(x int, y int) *Letter {
	return &lm.letters[x][y]
}

func (lm *LetterMap) Print() {
	fmt.Println("0 1 2 3 4 5 6 7 8 9 10")
	fmt.Println("----------------------")

	for i := 0; i < lm.Len(); i++ {
		line := lm.letters[i]

		letters := []string{}
		for _, el := range line {
			letters = append(letters, string(el.L))
		}
		fmt.Println(i, strings.Join(letters, " "))
	}
}

func createMap(content string) LetterMap {
	letterMap := LetterMap{
		letters: [][]Letter{},
	}

	lines := strings.Split(content, "\n")

	temp := []Letter{}
	for i, line := range lines {
		for j, c := range strings.TrimSpace(line) {
			if isXmas(c) {
				letter := Letter{
					L:   c,
					Pos: [2]int{i, j},
				}

				temp = append(temp, letter)
			}
		}

		letterMap.letters = append(letterMap.letters, temp)
		temp = []Letter{}
	}

	return letterMap
}

func lookupPoints(l Letter, lMap LetterMap) [][2]int {
	horizontalLength := lMap.Len()
	verticalLength := len(lMap.letters[0])

	directions := [][2]int{
		{1, 1},   // Bottom-right
		{0, 1},   // Right
		{1, 0},   // Down
		{1, -1},  // Bottom-left
		{-1, 1},  // Top-right
		{-1, -1}, // Top-left
		{0, -1},  // Left
		{-1, 0},  // Up
	}

	validPaths := [][2]int{}
	for _, direction := range directions {
		newX := l.Pos[0] + direction[0]
		newY := l.Pos[1] + direction[1]

		if newX >= 0 && newX < horizontalLength && newY >= 0 && newY < verticalLength {
			validPaths = append(validPaths, [2]int{newX, newY})
		}
	}

	return validPaths
}

type LetterInStack struct {
	letter    Letter
	direction *[2]int
}

func countXMAS(l Letter, lMap LetterMap) int {
	if l.L != 'X' {
		return 0
	}

	count := 0

	stack := []LetterInStack{
		{
			letter:    l,
			direction: nil,
		},
	}

	validCheck := map[rune]rune{
		'X': 'M',
		'M': 'A',
		'A': 'S',
	}

	for len(stack) > 0 {
		currentLetter := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		var lookUp [][2]int
		if currentLetter.direction != nil {
			lookUp = [][2]int{{currentLetter.letter.Pos[0] + currentLetter.direction[0], currentLetter.letter.Pos[1] + currentLetter.direction[1]}}
		} else {
			lookUp = lookupPoints(l, lMap)
		}

		validPaths := [][2]int{}
		for _, direction := range lookUp {
			if direction[0] >= 0 && direction[0] < lMap.Len() && direction[1] >= 0 && direction[1] < len(lMap.letters) {
				validPaths = append(validPaths, [2]int{direction[0], direction[1]})
			}
		}

		fmt.Println(validPaths)

		for _, point := range validPaths {
			lookedUpLetter := lMap.Get(point[0], point[1])

			if lookedUpLetter.L == validCheck[currentLetter.letter.L] {
				stack = append(stack, LetterInStack{
					letter:    *lookedUpLetter,
					direction: &[2]int{lookedUpLetter.Pos[0] - currentLetter.letter.Pos[0], lookedUpLetter.Pos[1] - currentLetter.letter.Pos[1]},
				})

				if lookedUpLetter.L == 'S' {
					count++
				}
			}
		}
	}

	return count
}

func countCrossMAS(l Letter, lMap LetterMap) int {
	if l.L != 'A' {
		return 0
	}

	count := 0

	if isValidCross(l, lMap) {
		count++
	}

	return count
}

func isValidCross(l Letter, lMap LetterMap) bool {
	lookUp := lookupPoints(l, lMap)
	if len(lookUp) < 8 {
		return false
	}

	x, y := l.Pos[0], l.Pos[1]

	// Check diagonals
	topLeft := lMap.Get(x-1, y-1)
	topRight := lMap.Get(x-1, y+1)
	bottomLeft := lMap.Get(x+1, y-1)
	bottomRight := lMap.Get(x+1, y+1)

	validTopLeft := (topLeft.L == 'M' && bottomRight.L == 'S') || (topLeft.L == 'S' && bottomRight.L == 'M')
	validTopRight := (topRight.L == 'M' && bottomLeft.L == 'S') || (topRight.L == 'S' && bottomLeft.L == 'M')

	return validTopLeft && validTopRight
}

func main() {
	content, _ := os.ReadFile("4_data.txt")

	letterMap := createMap(string(content))
	count := 0

	letterMap.Print()
	fmt.Println()

	for i := 0; i < letterMap.Len(); i++ {
		line := letterMap.letters[i]

		for j := 0; j < len(line); j++ {
			letter := *letterMap.Get(i, j)
			letterCount := countCrossMAS(letter, letterMap)

			count += letterCount
		}

		fmt.Println()
	}

	fmt.Println(count)
}
