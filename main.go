package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const TEST_DATA = `....#.....
.........#
..........
..#.......
.......#..
..........
.#..^.....
........#.
#.........
......#...
`

type Guard struct {
	coords    [2]int
	direction [2]int
}

func (g *Guard) GetSymbol() string {
	if g.direction[0] == 1 {
		return "v"
	} else if g.direction[0] == -1 {
		return "^"
	} else if g.direction[1] == 1 {
		return ">"
	} else if g.direction[1] == -1 {
		return "<"
	} else {
		return "."
	}
}

func NewGuard(coords [2]int) Guard {
	return Guard{
		coords:    coords,
		direction: [2]int{-1, 0},
	}
}

type Obstruction struct {
	coords [2]int
}

type Emulation struct {
	xLen int
	yLen int

	guard        Guard
	obstructions []Obstruction

	visitedCoords [][2]int
}

func (e *Emulation) move() {
	e.visitedCoords = append(e.visitedCoords, e.guard.coords)

	e.guard.coords[0] += e.guard.direction[0]
	e.guard.coords[1] += e.guard.direction[1]
}

func (e *Emulation) rotate() {
	if e.guard.direction[0] == 1 {
		e.guard.direction[0] = 0
		e.guard.direction[1] = -1
	} else if e.guard.direction[0] == -1 {
		e.guard.direction[0] = 0
		e.guard.direction[1] = 1
	} else if e.guard.direction[1] == 1 {
		e.guard.direction[0] = 1
		e.guard.direction[1] = 0
	} else if e.guard.direction[1] == -1 {
		e.guard.direction[0] = -1
		e.guard.direction[1] = 0
	}
}

func (e *Emulation) HasObstruction(x int, y int) bool {
	hasObstruction := false

	for _, obs := range e.obstructions {
		if obs.coords[0] == x && obs.coords[1] == y {
			return true
		}
	}

	return hasObstruction
}

func (e *Emulation) HasVisited(x int, y int) bool {
	hasVisited := false

	for _, visited := range e.visitedCoords {
		if visited[0] == x && visited[1] == y {
			return true
		}
	}

	return hasVisited
}

func (e *Emulation) Advance() {
	newX, newY := e.guard.coords[0]+e.guard.direction[0], e.guard.coords[1]+e.guard.direction[1]

	if e.HasObstruction(newX, newY) {
		e.rotate()
	} else {
		e.move()
	}
}

func (e *Emulation) CountUniqueVisited() int {
	cache := map[string]string{}

	for _, point := range e.visitedCoords {
		key := strconv.Itoa(point[0]) + ":" + strconv.Itoa(point[1])

		_, ok := cache[key]
		if !ok {
			cache[key] = key
		}
	}

	return len(cache) - 1
}

func (e *Emulation) HasLeftMap() bool {
	x, y := e.guard.coords[0], e.guard.coords[1]
	return x < 0 || y < 0 || x > e.xLen || y > e.yLen
}

func (e *Emulation) GetSymbol(x int, y int) string {
	if e.guard.coords[0] == x && e.guard.coords[1] == y {
		return e.guard.GetSymbol()
	} else if e.HasObstruction(x, y) {
		return "#"
	} else if e.HasVisited(x, y) {
		return "X"
	} else {
		return "."
	}
}

func (e *Emulation) Show() {
	result := ""
	for i := 0; i < e.xLen; i++ {
		temp := ""
		for j := 0; j < e.yLen; j++ {
			symbol := e.GetSymbol(i, j)
			temp += symbol
		}

		result += temp + "\n"
	}

	fmt.Print(result)
}

func NewEmulation(content string) Emulation {
	var guard Guard
	obstructions := []Obstruction{}

	var xLen int
	lines := strings.Split(content, "\n")

	for i, line := range lines {

		if xLen == 0 {
			xLen = len(line)
		}

		for j, char := range line {
			if char == '#' {
				obstructions = append(obstructions, Obstruction{
					coords: [2]int{i, j},
				})
			} else if char == '^' {
				guard = NewGuard([2]int{i, j})
			}
		}
	}

	return Emulation{
		xLen: xLen,
		yLen: len(lines),

		visitedCoords: [][2]int{},

		guard:        guard,
		obstructions: obstructions,
	}
}

func main() {
	data, _ := os.ReadFile("data.txt")
	emulation := NewEmulation(string(data))

	fmt.Println(emulation.xLen, emulation.yLen)

	for !emulation.HasLeftMap() {
		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()

		emulation.Show()
		emulation.Advance()
	}

	fmt.Println(emulation.CountUniqueVisited())
}
