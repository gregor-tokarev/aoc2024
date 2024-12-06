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

type Visited struct {
	coords    [2]int
	direction [2]int
	count     int
}

type Emulation struct {
	xLen int
	yLen int

	guardCursor  Guard
	guard        Guard
	obstructions []Obstruction

	visitedCoords []Visited
}

func (e *Emulation) move() {
	visited := Visited{
		direction: [2]int{e.guardCursor.direction[0], e.guardCursor.direction[1]},
		coords:    [2]int{e.guardCursor.coords[0], e.guardCursor.coords[1]},
		count:     0,
	}

	e.visitedCoords = append(e.visitedCoords, visited)

	e.guardCursor.coords[0] += e.guardCursor.direction[0]
	e.guardCursor.coords[1] += e.guardCursor.direction[1]
}

func (e *Emulation) rotate() {
	if e.guardCursor.direction[0] == 1 {
		e.guardCursor.direction[0] = 0
		e.guardCursor.direction[1] = -1
	} else if e.guardCursor.direction[0] == -1 {
		e.guardCursor.direction[0] = 0
		e.guardCursor.direction[1] = 1
	} else if e.guardCursor.direction[1] == 1 {
		e.guardCursor.direction[0] = 1
		e.guardCursor.direction[1] = 0
	} else if e.guardCursor.direction[1] == -1 {
		e.guardCursor.direction[0] = -1
		e.guardCursor.direction[1] = 0
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

func (e *Emulation) SetObstruction(x int, y int) {
	if !e.HasObstruction(x, y) {
		e.obstructions = append(e.obstructions, Obstruction{
			coords: [2]int{x, y},
		})
	}
}

func (e *Emulation) RemoveObstruction(x int, y int) {
	for i, obst := range e.obstructions {
		if obst.coords[0] == x && obst.coords[1] == y {
			e.obstructions = append(e.obstructions[:i], e.obstructions[i+1:]...)
			break
		}
	}
}

func (e *Emulation) HasVisited(x int, y int) bool {
	for _, visited := range e.visitedCoords {
		if visited.coords[0] == x && visited.coords[1] == y {
			return true
		}
	}

	return false
}

func (e *Emulation) BeenInSameDirection() bool {
	x, y := e.guardCursor.coords[0], e.guardCursor.coords[1]
	xDirection, yDirection := e.guardCursor.direction[0], e.guardCursor.direction[1]

	for _, visited := range e.visitedCoords {
		if visited.coords[0] == x && visited.coords[1] == y && visited.direction[0] == xDirection && visited.direction[1] == yDirection {
			return true
		}
	}

	return false
}

func (e *Emulation) Advance() {
	newX, newY := e.guardCursor.coords[0]+e.guardCursor.direction[0], e.guardCursor.coords[1]+e.guardCursor.direction[1]

	if e.HasObstruction(newX, newY) {
		e.rotate()
	} else {
		e.move()
	}
}

func (e *Emulation) CountUniqueVisited() int {
	cache := map[string]string{}

	for _, point := range e.visitedCoords {
		key := strconv.Itoa(point.coords[0]) + ":" + strconv.Itoa(point.coords[1])

		_, ok := cache[key]
		if !ok {
			cache[key] = key
		}
	}

	return len(cache) - 1
}

func (e *Emulation) HasLeftMap() bool {
	x, y := e.guardCursor.coords[0], e.guardCursor.coords[1]
	return x < 0 || y < 0 || x > e.xLen || y > e.yLen
}

func (e *Emulation) GetSymbol(x int, y int) string {
	if e.guardCursor.coords[0] == x && e.guardCursor.coords[1] == y {
		return e.guardCursor.GetSymbol()
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

	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()

	fmt.Print(result)
}

func (e *Emulation) Reset() {
	e.guardCursor = NewGuard([2]int{e.guard.coords[0], e.guard.coords[1]})
	e.visitedCoords = []Visited{}
}

func (e *Emulation) HasLoop() bool {
	for !e.HasLeftMap() {
		e.Advance()

		if e.BeenInSameDirection() {
			return true
		}
	}

	return false
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

		visitedCoords: []Visited{},

		guardCursor:  guard,
		guard:        NewGuard([2]int{guard.coords[0], guard.coords[1]}),
		obstructions: obstructions,
	}
}

func main() {
	data, _ := os.ReadFile("data.txt")
	emulation := NewEmulation(string(data))

	count := 0

	for i := 0; i < emulation.xLen; i++ {
		for j := 0; j < emulation.yLen; j++ {
			if emulation.HasObstruction(i, j) {
				continue
			}

			emulation.SetObstruction(i, j)

			if emulation.HasLoop() {
				count++
			}

			emulation.RemoveObstruction(i, j)
			emulation.Reset()
		}
	}

	fmt.Println(count)
}
