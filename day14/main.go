package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func getNums(block string) [2]int {
	splitted := strings.Split(block, "=")

	nums := strings.Split(splitted[1], ",")

	res := []int{}
	for _, ns := range nums {
		n, _ := strconv.Atoi(ns)
		res = append(res, n)
	}

	return [2]int(res)
}

type Arena struct {
	robots []Robot

	xLen int
	yLen int
}

func (a *Arena) GetRobots(x int, y int) []Robot {
	robots := []Robot{}

	for _, r := range a.robots {
		if r.pos[0] == x && r.pos[1] == y {
			robots = append(robots, r)
		}
	}

	return robots
}

func (a *Arena) Show() {
	var res strings.Builder

	for y := 0; y < a.yLen; y++ {
		for x := 0; x < a.xLen; x++ {
			robots := a.GetRobots(x, y)
			if len(robots) > 0 {
				res.WriteString(strconv.Itoa(len(robots)))
			} else {
				res.WriteString(".")
			}
		}

		res.WriteString("\n")
	}

	fmt.Print(res.String() + "\n")
}

func FromStrArena(content string) Arena {
	robots := []Robot{}

	for _, line := range strings.Split(content, "\n") {
		robots = append(robots, FromStrRobot(line))
	}

	return Arena{
		robots: robots,

		xLen: 101,
		yLen: 103,
	}
}

func (a *Arena) Advance() {
	for idx := range a.robots {
		r := &a.robots[idx]

		r.pos[0] += r.direction[0]
		r.pos[1] += r.direction[1]

		if r.pos[0] >= a.xLen {
			r.pos[0] = r.pos[0] - a.xLen
		} else if r.pos[0] < 0 {
			r.pos[0] = a.xLen + r.pos[0]
		}

		if r.pos[1] >= a.yLen {
			r.pos[1] = r.pos[1] - a.yLen
		} else if r.pos[1] < 0 {
			r.pos[1] = a.yLen + r.pos[1]
		}
	}
}

func (a *Arena) GetProduct() int {
	counts := [4]int{}

	xBorder, yBorder := a.xLen/2, a.yLen/2
	for _, r := range a.robots {
		if r.pos[0] > xBorder && r.pos[1] < yBorder {
			counts[0]++
		}

		if r.pos[0] > xBorder && r.pos[1] > yBorder {
			counts[1]++
		}

		if r.pos[0] < xBorder && r.pos[1] > yBorder {
			counts[2]++
		}

		if r.pos[0] < xBorder && r.pos[1] < yBorder {
			counts[3]++
		}
	}

	sum := 1
	for _, count := range counts {
		sum *= count
	}

	return sum
}

type Robot struct {
	pos       [2]int
	direction [2]int
}

func FromStrRobot(line string) Robot {
	splitted := strings.Split(line, " ")

	pos := getNums(splitted[0])
	direction := getNums(splitted[1])

	return Robot{
		pos:       pos,
		direction: direction,
	}
}

const test_data = `p=0,4 v=3,-3
p=6,3 v=-1,-3
p=10,3 v=-1,2
p=2,0 v=2,-1
p=0,0 v=1,3
p=3,0 v=-2,-2
p=7,6 v=-1,-3
p=3,0 v=-1,-2
p=9,3 v=2,3
p=7,3 v=-1,2
p=2,4 v=2,-3
p=9,5 v=-3,-3`

const test_data2 = `p=2,4 v=2,-3`

func main() {
	data, _ := os.ReadFile("data.txt")
	arena := FromStrArena(string(data))

	for i := 0; i < 100; i++ {
		arena.Advance()

		if i > 50 {
			fmt.Println(i)
			arena.Show()

			time.Sleep(100000 * time.Microsecond)
		}
	}

	fmt.Println(arena.GetProduct())
}
