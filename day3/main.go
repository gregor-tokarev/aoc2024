package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MulInstraction struct {
	left  string
	right string
}

func (i *MulInstraction) mulResult() (int, error) {
	leftNum, err := strconv.Atoi(i.left)
	if err != nil {
		return 0, err
	}

	rightNum, err := strconv.Atoi(i.right)
	if err != nil {
		return 0, err
	}

	return leftNum * rightNum, nil
}

func isInteger(check string) bool {
	_, err := strconv.Atoi(check)
	return err == nil
}

var brackets = map[string]string{
	"(": ")",
}

func scrapeMuls(content string) []MulInstraction {
	instructions := []MulInstraction{}

	accepting := true

	i := 0
	for i < len(content) {
		if content[i] == 'd' && content[i:i+4] == "do()" {
			accepting = true

			i += 4
			continue
		} else if content[i] == 'd' && content[i:i+7] == "don't()" {
			accepting = false

			i += 7
			continue
		} else if content[i] == 'm' && content[i:i+3] == "mul" {
			i += 3

			if content[i:i+1] != "[" && content[i:i+1] != "(" {
				continue
			}
			bracket := content[i : i+1]
			i++

			instruction := MulInstraction{
				left:  "",
				right: "",
			}

			var currentStringBuilder strings.Builder

			for isInteger(string(content[i])) || content[i:i+1] == "," {
				if content[i] == ',' {
					instruction.left = currentStringBuilder.String()

					currentStringBuilder.Reset()
				} else if isInteger(string(content[i])) {
					currentStringBuilder.WriteByte(content[i])
				}

				i++
			}

			instruction.right = currentStringBuilder.String()

			pair, ok := brackets[bracket]

			if !ok || pair != content[i:i+1] || !accepting {
				continue
			}

			instructions = append(instructions, instruction)
		}

		i++
	}

	return instructions
}

func main() {
	content, _ := os.ReadFile("3_data.txt")
	// content := "xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))"

	instructions := scrapeMuls(string(content))

	sum := 0
	for _, inst := range instructions {
		result, _ := inst.mulResult()
		sum += result
	}

	fmt.Println(sum)
}
