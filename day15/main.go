package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const data = `########
#..O.O.#
##@.O..#
#...O..#
#.#.O..#
#...O..#
#......#
########

<^^>>>vv<v>>v<<`

type Move int

func (m Move) GetDirection() [2]int {
	switch m {
	case Top:
		return [2]int{0, -1}
	case Right:
		return [2]int{1, 0}
	case Bottom:
		return [2]int{0, 1}
	case Left:
		return [2]int{-1, 0}
	default:
		return [2]int{0, 0}
	}
}

const (
	Top Move = iota
	Right
	Bottom
	Left
)

type Box struct {
	x, y int
}

type Wall struct {
	x, y int
}

type Robot struct {
	pos [2]int
}

type Arena struct {
	xLen int
	yLen int

	moves []Move

	boxes []Box
	walls []Wall

	robot Robot
}

func Clear() {
	cmd := exec.Command("clear") // Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (a *Arena) Show() {
	var builder strings.Builder

	for y := 0; y < a.yLen; y++ {
		for x := 0; x < a.xLen; x++ {
			if a.GetWall(x, y) != nil {
				builder.WriteRune('#')
			} else if a.GetBox(x, y) != nil {
				builder.WriteRune('O')
			} else if a.robot.pos[0] == x && a.robot.pos[1] == y {
				builder.WriteRune('@')
			} else {
				builder.WriteRune('.')
			}
		}
		builder.WriteRune('\n')
	}

	fmt.Print(builder.String())
}

func (a *Arena) GetWall(x int, y int) *Wall {
	var wall *Wall

	for _, w := range a.walls {
		if w.x == x && w.y == y {
			wall = &w
			break
		}
	}

	return wall
}

func (a *Arena) GetBox(x int, y int) *Box {
	for i := range a.boxes {
		if a.boxes[i].x == x && a.boxes[i].y == y {
			return &a.boxes[i]
		}
	}
	return nil
}

func (a *Arena) GetBoxChain(pos [2]int, direction [2]int) []*Box {
	boxes := []*Box{}

	boxCursor := a.GetBox(pos[0], pos[1])
	for boxCursor != nil {
		boxes = append(boxes, boxCursor)

		newX, newY := boxCursor.x+direction[0], boxCursor.y+direction[1]
		boxCursor = a.GetBox(newX, newY)
	}

	return boxes
}

func (a *Arena) CountBoxesGps() int {
	sum := 0

	for _, box := range a.boxes {
		sum += box.x + 100*box.y
	}

	return sum
}

func (a *Arena) Emulate() {
	for _, move := range a.moves {
		direction := move.GetDirection()
		x, y := a.robot.pos[0]+direction[0], a.robot.pos[1]+direction[1]

		if a.GetWall(x, y) != nil {
			continue
		} else if a.GetBox(x, y) != nil {
			boxChain := a.GetBoxChain([2]int{x, y}, direction)

			lastBox := boxChain[len(boxChain)-1]

			if a.GetWall(lastBox.x+direction[0], lastBox.y+direction[1]) != nil {
				continue
			}

			for _, box := range boxChain {
				box.x += direction[0]
				box.y += direction[1]
			}

			a.robot.pos[0] = x
			a.robot.pos[1] = y
		} else {
			a.robot.pos[0] = x
			a.robot.pos[1] = y
		}
	}
}

func NewArena(content string) Arena {
	splitted := strings.Split(content, "\n\n")

	robotMapLines := strings.Split(splitted[0], "\n")
	movesLine := splitted[1]

	var robot Robot

	boxes := []Box{}
	walls := []Wall{}
	for y, line := range robotMapLines {
		for x, char := range line {
			switch char {
			case '#':
				walls = append(walls, Wall{x, y})
			case 'O':
				boxes = append(boxes, Box{x, y})
			case '@':
				robot = Robot{pos: [2]int{x, y}}
			}
		}
	}

	moves := []Move{}
	for _, move := range movesLine {
		switch move {
		case '^':
			moves = append(moves, Top)
		case '>':
			moves = append(moves, Right)
		case 'v':
			moves = append(moves, Bottom)
		case '<':
			moves = append(moves, Left)
		}
	}

	return Arena{
		xLen: len(robotMapLines[0]),
		yLen: len(robotMapLines),

		moves: moves,

		boxes: boxes,
		walls: walls,

		robot: robot,
	}
}

func main() {
	puzzle, _ := os.ReadFile("data.txt")
	arena := NewArena(string(puzzle))
	arena.Emulate()

	arena.Show()

	fmt.Println(arena.CountBoxesGps())
}
