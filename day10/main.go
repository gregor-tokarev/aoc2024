package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const test_data = `89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732`

type TMap struct {
	heights [][]int

	xLen int
	yLen int
}

func NewTMap(content string) TMap {
	heights := [][]int{}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		tmp := []int{}
		for _, char := range line {
			num, _ := strconv.Atoi(string(char))
			tmp = append(tmp, num)
		}

		heights = append(heights, tmp)
	}

	return TMap{
		heights: heights,

		yLen: len(lines),
		xLen: len(lines[0]),
	}
}

var directions = [][2]int{
	{0, 1},
	{1, 0},
	{0, -1},
	{-1, 0},
}

func (m *TMap) calculateScore(x int, y int) int {
	stack := [][2]int{{y, x}}
	visited := [][2]int{}

	endPositions := [][2]int{}

	for len(stack) > 0 {
		coords := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		currentHeight := m.heights[coords[0]][coords[1]]
		// fmt.Println(currentHeight, coords, visited)

		if currentHeight == 9 {
			hasSamePosition := false

			for _, pos := range endPositions {
				if pos[0] == coords[0] && pos[1] == coords[1] {
					hasSamePosition = true
				}
			}

			if !hasSamePosition {
				endPositions = append(endPositions, [2]int{coords[0], coords[1]})
			}

			continue
		}

		for _, dir := range directions {
			nextY := coords[0] + dir[0]
			nextX := coords[1] + dir[1]

			if nextX >= m.xLen || nextY >= m.yLen || nextX < 0 || nextY < 0 {
				continue
			}

			next := m.heights[nextY][nextX]

			if next == currentHeight+1 {
				stack = append(stack, [2]int{nextY, nextX})
			}
		}

		visited = append(visited, coords)
	}

	return len(endPositions)
}

func main() {
	data, _ := os.ReadFile("data.txt")
	tMap := NewTMap(string(data))

	sum := 0
	for y := 0; y < tMap.yLen; y++ {
		for x := 0; x < tMap.xLen; x++ {
			if tMap.heights[y][x] == 0 {
				sum += tMap.calculateScore(x, y)
			}
		}
	}

	fmt.Println(sum)
}
