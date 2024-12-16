package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const TEST_DATA = "2333133121414131402"

type MemorySlot struct {
	count  int
	ID     *int
	isFree bool
}

type Disk struct {
	slots []MemorySlot

	freeMemoryCursor int
	rightCursor      int
}

func NewDisk(data string) Disk {
	slots := []MemorySlot{}

	for i, char := range data {
		num, _ := strconv.Atoi(string(char))

		if i%2 != 0 {
			slots = append(slots, MemorySlot{
				count:  num,
				ID:     nil,
				isFree: true,
			})
		} else {
			var id int
			if i == 0 {
				id = 0
			} else {
				id = i / 2
			}

			slots = append(slots, MemorySlot{
				count:  num,
				ID:     &id,
				isFree: false,
			})
		}
	}

	var rightCursor int
	if len(slots)%2 != 0 {
		rightCursor = len(slots) - 1
	} else {
		rightCursor = len(slots) - 2
	}

	return Disk{
		slots: slots,

		freeMemoryCursor: 1,
		rightCursor:      rightCursor,
	}
}

func (d *Disk) Show() {
	var result strings.Builder

	for _, slot := range d.slots {

		var c string
		if slot.isFree {
			c = "."
		} else {
			if slot.ID != nil {
				c = strconv.Itoa(*slot.ID)
			}
		}

		result.WriteString(strings.Repeat(c, slot.count))
	}

	fmt.Println(result.String())
}

func (d *Disk) Advance() {
	freeMemory := d.slots[d.freeMemoryCursor]
	block := d.slots[d.rightCursor]

	if freeMemory.count <= 0 {
		d.freeMemoryCursor += 2
		return
	} else if block.count <= 0 {
		d.rightCursor -= 2
		return
	}

	if freeMemory.count > block.count {
		newSlot := MemorySlot{
			count:  block.count,
			ID:     block.ID,
			isFree: false,
		}
		d.slots = append(d.slots[:d.freeMemoryCursor], append([]MemorySlot{newSlot}, d.slots[d.freeMemoryCursor:]...)...)

		// Update counts
		d.slots[d.freeMemoryCursor+1].count = freeMemory.count - block.count
		d.slots[d.rightCursor].count = 0

		d.rightCursor -= 2
	} else if freeMemory.count < block.count {
		newSlot := MemorySlot{
			count:  freeMemory.count,
			ID:     block.ID,
			isFree: false,
		}
		d.slots = append(d.slots[:d.freeMemoryCursor], append([]MemorySlot{newSlot}, d.slots[d.freeMemoryCursor:]...)...)

		// Update counts
		d.slots[d.rightCursor].count -= freeMemory.count
		d.slots[d.freeMemoryCursor+1].count = 0

		d.freeMemoryCursor += 2
	} else {
		d.rightCursor -= 2
		d.freeMemoryCursor += 2
	}
}

func main() {
	disk := NewDisk(TEST_DATA)

	for disk.freeMemoryCursor < disk.rightCursor {
		disk.Show()
		disk.Advance()

		time.Sleep(10 * time.Millisecond)
	}
}
