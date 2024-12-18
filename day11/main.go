package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Stone struct {
	num string
}

func (s *Stone) GetInt() int {
	num, _ := strconv.Atoi(s.num)

	return num
}

func (s Stone) SplitParts() (Stone, Stone) {
	leftPart, rightPart := s.num[:len(s.num)/2], s.num[len(s.num)/2:]

	return Stone{num: leftPart}, Stone{num: rightPart}
}

func (s Stone) Multiply() Stone {
	newNum := s.GetInt() * 2024

	return Stone{
		num: strconv.Itoa(newNum),
	}
}

type StoneSeq struct {
	stones []Stone
	cache  map[string][]Stone
}

func (ss *StoneSeq) Advance() {
	idx := 0
	maxIdx := len(ss.stones)

	for idx < maxIdx {
		stone := ss.stones[idx]
		stone.num = strings.TrimLeft(stone.num, "0")

		if stone.num == strings.Repeat("0", len(stone.num)) {
			ss.stones[idx] = Stone{num: "1"}
		} else if len(stone.num)%2 != 0 {
			ss.stones[idx] = stone.Multiply()
		} else {
			stone1, stone2 := stone.SplitParts()
			ss.stones = slices.Replace(ss.stones, idx, idx+1, stone1, stone2)

			idx++
			maxIdx++
		}

		idx++
	}
}

func (ss *StoneSeq) AdvanceEfficient() {
	output := make([]Stone, 0, len(ss.stones)*2)

	for _, stone := range ss.stones {
		stone.num = strings.TrimLeft(stone.num, "0")

		if stone.num == strings.Repeat("0", len(stone.num)) {
			output = append(output, Stone{num: "1"})
		} else if len(stone.num)%2 != 0 {
			output = append(output, stone.Multiply())
		} else {
			stone1, stone2 := stone.SplitParts()
			output = append(output, stone1, stone2)
		}
	}

	ss.stones = output
}

func (ss *StoneSeq) AdvanceWithCache() { // Takes more then 30gig of memory on my mac
	output := make([]Stone, 0, len(ss.stones)*2)

	for _, stone := range ss.stones {
		stone.num = strings.TrimLeft(stone.num, "0")

		// Check cache first
		if val, ok := ss.cache[stone.num]; ok {
			output = append(output, val...)
			continue
		}

		// Process stone and cache the result
		var result []Stone
		if stone.num == "" || stone.num == "0" { // Handle empty or zero strings
			result = []Stone{{num: "1"}}
		} else if len(stone.num)%2 != 0 {
			multiplied := stone.Multiply()
			result = []Stone{multiplied}
		} else {
			stone1, stone2 := stone.SplitParts()
			result = []Stone{stone1, stone2}
		}

		// Update cache and output
		ss.cache[stone.num] = result
		output = append(output, result...)
	}

	ss.stones = output
}

func NewStoneSeq(content string) StoneSeq {
	stones := []Stone{}

	for _, num := range strings.Split(content, " ") {
		stones = append(stones, Stone{num: num})
	}

	return StoneSeq{
		stones: stones,
		cache:  map[string][]Stone{},
	}
}

func main() {
	stoneSeq := NewStoneSeq("4022724 951333 0 21633 5857 97 702 6")
	// stoneSeq := NewStoneSeq("125 17")

	for i := 0; i < 75; i++ {
		timeBefore := time.Now()
		stoneSeq.AdvanceEfficient()
		fmt.Println(i, time.Since(timeBefore))
	}

	fmt.Println(len(stoneSeq.stones))
}
