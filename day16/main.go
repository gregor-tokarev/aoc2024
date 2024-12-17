package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const data = `###############
#.......#....E#
#.#.###.#.###.#
#.....#.#...#.#
#.###.#####.#.#
#.#.#.......#.#
#.#.#####.###.#
#...........#.#
###.#.#####.#.#
#...#.....#.#.#
#.#.#.###.#.#.#
#.....#...#.#.#
#.###.#.#.#.#.#
#S..#.....#...#
###############`

type Wall struct {
	x, y int
}

type Cursor struct {
	direction [2]int
	pos       [2]int
}

func NewCursor(pos [2]int) Cursor {
	return Cursor{
		direction: [2]int{0, 1},
		pos:       pos,
	}
}

func (c *Cursor) move() {
	c.pos[0] += c.direction[0]
	c.pos[1] += c.direction[1]
}

type RotateDirection int

const (
	Right RotateDirection = iota
	Left
)

func (c Cursor) rotate(direction RotateDirection) Cursor {
	var newDirection [2]int
	switch direction {
	case Right:
		// Rotate clockwise: (x, y) -> (y, -x)
		newDirection = [2]int{c.direction[1], -c.direction[0]}
	case Left:
		// Rotate counter-clockwise: (x, y) -> (-y, x)
		newDirection = [2]int{-c.direction[1], c.direction[0]}
	}

	// Return a new Cursor with the updated direction, keeping position the same
	return Cursor{
		direction: newDirection,
		pos:       c.pos,
	}
}

type Game struct {
	walls []Wall

	start [2]int
	end   [2]int
}

func (g *Game) isGameEnd(cursor Cursor) bool {
	return cursor.pos[0] == g.end[0] && cursor.pos[1] == g.end[1]
}

func (g *Game) findShortestPath() int {
	type State struct {
		pos       [2]int
		direction [2]int
		score     int
	}
	directions := [][2]int{
		{0, 1},  // East
		{1, 0},  // South
		{0, -1}, // West
		{-1, 0}, // North
	}

	visited := make(map[[3]int]bool)             // Key: [x, y, directionIndex]
	queue := []State{{g.start, [2]int{0, 1}, 0}} // Start facing East

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// Check if we reached the end
		if current.pos == g.end {
			return current.score
		}

		// Encode the state to prevent revisiting
		dirIndex := -1
		for i, d := range directions {
			if d == current.direction {
				dirIndex = i
				break
			}
		}
		key := [3]int{current.pos[0], current.pos[1], dirIndex}
		if visited[key] {
			continue
		}
		visited[key] = true

		// Explore forward move
		nextPos := [2]int{current.pos[0] + current.direction[0], current.pos[1] + current.direction[1]}
		if !g.hasWall(nextPos[0], nextPos[1]) {
			queue = append(queue, State{nextPos, current.direction, current.score + 1})
		}

		// Explore rotations
		for _, dir := range directions {
			if dir == current.direction {
				newScore := current.score + 1
				queue = append(queue, State{current.pos, dir, newScore})
				continue
			}
			newScore := current.score + 1000
			queue = append(queue, State{current.pos, dir, newScore})
		}
	}

	return -1 // No path found
}

func (g *Game) isFacingDeadEnd(cursor Cursor) bool {
	coords := [][2]int{
		{cursor.pos[0] + 1, cursor.pos[1]}, // Right
		{cursor.pos[0] - 1, cursor.pos[1]}, // Left
		{cursor.pos[0], cursor.pos[1] + 1}, // Down
		{cursor.pos[0], cursor.pos[1] - 1},
	}

	wallCount := 0
	for _, coord := range coords {
		if g.hasWall(coord[0], coord[1]) {
			wallCount++
		}
	}

	return wallCount >= 3
}

func (g *Game) hasWall(x int, y int) bool {
	var wall *Wall
	for _, w := range g.walls {
		if w.x == x && w.y == y {
			wall = &w
		}
	}

	return wall != nil
}

func (g *Game) isFacingWall(cursor Cursor) bool {
	x := cursor.pos[0] + cursor.direction[0]
	y := cursor.pos[1] + cursor.direction[1]

	return g.hasWall(x, y)
}

func NewGame(data string) Game {
	lines := strings.Split(data, "\n")

	walls := []Wall{}
	var start, end [2]int

	for y, line := range lines {
		for x, itm := range line {
			if itm == '#' {
				walls = append(walls, Wall{x, y})
			} else if itm == 'S' {
				start = [2]int{x, y}
			} else if itm == 'E' {
				end = [2]int{x, y}
			}
		}
	}

	return Game{
		walls: walls,

		start: start,
		end:   end,
	}
}

func (g *Game) emulate(score int, cursor Cursor) int {
	fmt.Println(score, cursor)
	if g.isFacingDeadEnd(cursor) {
		return int(math.Inf(1))
	} else if g.isGameEnd(cursor) {
		return score
	}

	time.Sleep(10000 * time.Microsecond)

	results := []int{}

	if !g.isFacingWall(cursor) {
		newCursor := Cursor{
			pos:       [2]int{cursor.pos[0], cursor.pos[1]},
			direction: [2]int{cursor.direction[0], cursor.direction[1]},
		}
		newCursor.move()
		results = append(results, g.emulate(score+1, newCursor))
	}

	directions := [][2]int{
		{0, -1},
		{-1, 0},
		{0, 1},
		{1, 0},
	}

	for _, dir := range directions {
		newCursor := Cursor{
			pos:       [2]int{cursor.pos[0] + dir[0], cursor.pos[1] + dir[1]},
			direction: dir,
		}

		results = append(results, g.emulate(score+1000, newCursor))
	}

	minScore := int(math.Inf(1))
	for _, sc := range results {
		if minScore > sc {
			minScore = sc
		}
	}

	return minScore
}

func main() {
	game := NewGame(data)

	score := game.findShortestPath()
	fmt.Println(score)
}
