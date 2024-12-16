package main

import (
	"fmt"
	"os"
	"strings"
)

const test_data = `AAAAAA
AAABBA
AAABBA
ABBAAA
ABBAAA
AAAAAA`

func DEFAULT_DIRECTIONS() [4][2]int {
	return [4][2]int{
		{-1, 0},
		{0, -1},
		{1, 0},
		{0, 1},
	}
}

type Plot struct {
	x int
	y int

	name rune
}

type Garden struct {
	garenMap [][]Plot // x * y map
	visited  [][2]int

	xLen int
	yLen int

	areas [][]*Plot // slice of areas
}

func (g *Garden) GetPlot(x int, y int) *Plot {
	return &g.garenMap[y][x]
}

func (g *Garden) isVisited(x int, y int) bool {
	for _, coords := range g.visited {
		if coords[0] == x && coords[1] == y {
			return true
		}
	}

	return false
}

func (g *Garden) calculatePriceWithPerimetr() int {
	result := 0

	for _, area := range g.areas {
		perimetr := 0
		for _, point := range area {
			neiborsCount := len(g.getNeibors(point.x, point.y))
			perimetr += 4 - neiborsCount
		}

		result += perimetr * len(area)
	}

	return result
}

// Deprecated
func (g *Garden) calculatePriceWithSides() int {
	result := 0

	for _, area := range g.areas {
		sides := 0

		prevPlotNeiborsDrections := [][2]int{}

		for _, point := range area {
			fmt.Println(string(point.name))
			neibors := g.getNeibors(point.x, point.y)

			freeDirections := [][2]int{}
		outer:
			for _, d := range DEFAULT_DIRECTIONS() {
				neiborCoords := [2]int{point.x + d[0], point.y + d[1]}

				for _, n := range neibors {
					if n.x == neiborCoords[0] && n.y == neiborCoords[1] {
						continue outer
					}
				}

				freeDirections = append(freeDirections, d)
			}

			freeDirectionsWithoutPrevDirections := [][2]int{}
		outer2:
			for _, d := range freeDirections {
				for _, prevD := range prevPlotNeiborsDrections {
					if prevD[0] == d[0] && prevD[1] == d[1] {
						continue outer2
					}
				}

				freeDirectionsWithoutPrevDirections = append(freeDirectionsWithoutPrevDirections, d)
			}

			prevPlotNeiborsDrections = freeDirections
			fmt.Println(freeDirections, freeDirectionsWithoutPrevDirections)

			sides += len(freeDirectionsWithoutPrevDirections)
		}

		fmt.Println(sides)
		result += sides * len(area)
	}

	return result
}

func (g *Garden) getNeibors(x int, y int) []*Plot {
	basePlot := g.GetPlot(x, y)
	if basePlot == nil {
		return []*Plot{}
	}

	directions := DEFAULT_DIRECTIONS()
	neibors := []*Plot{}

	for _, direction := range directions {
		newX := x + direction[0]
		newY := y + direction[1]
		if newX < 0 || newX > (g.xLen-1) || newY < 0 || newY > (g.yLen-1) {
			continue
		}

		plot := g.GetPlot(newX, newY)

		if plot.name == basePlot.name {
			neibors = append(neibors, plot)
		}
	}

	return neibors
}

func NewGarden(content string) Garden {
	lines := strings.Split(content, "\n")

	garden := Garden{
		garenMap: [][]Plot{},
		visited:  [][2]int{},

		xLen: len(lines[0]),
		yLen: len(lines),

		areas: [][]*Plot{},
	}

	for y, line := range lines {
		tmp := []Plot{}

		for x, el := range line {
			tmp = append(tmp, Plot{x: x, y: y, name: el})
		}

		garden.garenMap = append(garden.garenMap, tmp)
		tmp = []Plot{}
	}

	for _, row := range garden.garenMap {
		for _, col := range row {
			stack := []*Plot{}
			stack = append(stack, &col)

			area := []*Plot{}

			for len(stack) > 0 {
				itm := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				if garden.isVisited(itm.x, itm.y) {
					continue
				}

				area = append(area, itm)
				garden.visited = append(garden.visited, [2]int{itm.x, itm.y})

				neibors := []*Plot{}
				for _, n := range garden.getNeibors(itm.x, itm.y) {
					if !garden.isVisited(n.x, n.y) {
						neibors = append(neibors, n)
					}
				}

				stack = append(stack, neibors...)
			}

			if len(area) > 0 {
				garden.areas = append(garden.areas, area)
				area = []*Plot{}
			}
		}
	}

	return garden
}

func main() {
	data, _ := os.ReadFile("data.txt")
	garden := NewGarden(string(data))

	fmt.Println(garden.calculatePriceWithPerimetr())
}
