package main

import (
	"fmt"
	"os"
	"strings"
)

const TEST_DATA = `T....#....
...T......
.T....#...
.........#
..#.......
..........
...#......
..........
....#.....
..........`

type Antenna struct {
	freq rune
	pos  [2]int
}

type Antinode struct {
	pos [2]int
}

type AntennaMap struct {
	antennas []Antenna

	lenX int
	lenY int

	antinodes []Antinode
}

func (am AntennaMap) Show() {
	visualMap := ""

	for i := 0; i < am.lenY; i++ {
		line := ""

		for j := 0; j < am.lenX; j++ {
			antenna := am.GetAnntena(j, i)
			if antenna != nil {
				line += string(antenna.freq)
			} else if am.GetAntinode(j, i) != nil {
				line += "#"
			} else {
				line += "."
			}
		}

		visualMap += line
		visualMap += "\n"
	}

	fmt.Print(visualMap)
}

func (am *AntennaMap) GetAnntena(x int, y int) *Antenna {
	var antenna *Antenna

	for _, a := range am.antennas {
		if a.pos[0] == x && a.pos[1] == y {
			antenna = &a
		}
	}

	return antenna
}

func (am *AntennaMap) GetSameFreq(x int, y int) []Antenna {
	sameFreq := []Antenna{}

	targetAntenna := am.GetAnntena(x, y)
	if targetAntenna == nil {
		return sameFreq
	}

	for _, antenna := range am.antennas {
		if antenna.freq != targetAntenna.freq || (antenna.pos[0] == x && antenna.pos[1] == y) {
			continue
		}

		sameFreq = append(sameFreq, antenna)
	}

	return sameFreq
}

func (am *AntennaMap) GetAntinode(x int, y int) *Antinode {
	var antinode *Antinode
	for _, anti := range am.antinodes {
		if anti.pos[0] == x && anti.pos[1] == y {
			antinode = &anti
		}
	}

	return antinode
}

func (am *AntennaMap) SetAntinode(x int, y int) {
	if am.lenX <= x || am.lenY <= y || x < 0 || y < 0 {
		return
	}

	if am.GetAntinode(x, y) == nil {
		am.antinodes = append(am.antinodes, Antinode{
			pos: [2]int{x, y},
		})
	}
}

func (am *AntennaMap) SetAntinodes() {
	for _, antenna := range am.antennas {
		sameFreq := am.GetSameFreq(antenna.pos[0], antenna.pos[1])

		for _, freqEntry := range sameFreq {
			diffX, diffY := freqEntry.pos[0]-antenna.pos[0], freqEntry.pos[1]-antenna.pos[1]

			x, y := antenna.pos[0], antenna.pos[1]

			for x < am.lenX && y < am.lenY && x >= 0 && y >= 0 {
				newX := x + diffX
				newY := y + diffY

				am.SetAntinode(newX, newY)

				x = newX
				y = newY
			}
		}
	}
}

func parseInput(data string) AntennaMap {
	splitted := strings.Split(data, "\n")

	antennas := []Antenna{}

	for y, line := range splitted {
		for x, char := range line {
			if char == '.' || char == '#' {
				continue
			}

			antennas = append(antennas, Antenna{
				freq: char,
				pos:  [2]int{x, y},
			})
		}
	}

	return AntennaMap{
		antennas: antennas,

		lenX: len(splitted[0]),
		lenY: len(splitted),

		antinodes: []Antinode{},
	}
}

func main() {
	data, _ := os.ReadFile("data.txt")
	anntennaMap := parseInput(string(data))

	anntennaMap.SetAntinodes()
	anntennaMap.Show()

	fmt.Println(len(anntennaMap.antinodes))
}
