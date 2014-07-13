package conway

import (
	"log"
	"bytes"
	"strconv"
)

type Cell int

const (
	empty = 0
	alive = 1
	//	dead  = 2
)

type Board struct {
	cells [][]Cell
}

func (this* Board) Init(size int) {
	this.cells = make([][]Cell, size)
	for i := range this.cells {
		this.cells[i] = make([]Cell, size)
	}

	log.Printf("Initialized board with size: [%d][%d]", len(this.cells), len(this.cells[0]))
}

func (this* Board) InitWithSeed(size int, seed [][]Cell) {
	this.Init(size)
	for x, row := range seed {
		for y, cell := range row {
			this.cells[x][y] = cell
		}
	}
	log.Printf("Initialized board with size and seed: [%d][%d] %v", len(this.cells), len(this.cells[0]), seed)
}

func (this* Board) Reaper() {
	log.Printf("Running reaper")
	board := make([][]Cell, len(this.cells))
	copy(board, this.cells)
	for x, row := range this.cells {
		copy(board[x], row)
	}

	for x, row := range this.cells {
		for y, _ := range row {
			neighbours := countNeighbours(x, y, board)

			cell := this.cell(x, y)
			if cell == empty {
				if neighbours == 3 {
					this.cells[x][y] = alive
				}
			} else if cell == alive {
				if neighbours == 2 || neighbours == 3 {
					this.cells[x][y] = alive
				} else {
					this.cells[x][y] = empty
				}
			}
		}
	}
}

func (this* Board) cell(x int, y int) Cell {
	return this.cells[x][y]
}

func countNeighbours(x int, y int, board [][]Cell) int {
	cell := board[x][y]
	if cell != empty {
		log.Printf("Cell: %d [%d, %d]", cell, x, y)
	}

	neighbours := 0
	if (x > 0 && y > 0) && (x < len(board)-1 && y < len(board)-1 ) {
		//Top row
		neighbours = neighbours+int(board[x - 1][y - 1])
		neighbours = neighbours+int(board[x][y - 1])
		neighbours = neighbours+int(board[x + 1][y - 1])
		//Middle
		neighbours = neighbours+int(board[x - 1][y])
		neighbours = neighbours+int(board[x + 1][y])
		//Bottom row
		neighbours = neighbours+int(board[x - 1][y + 1])
		neighbours = neighbours+int(board[x][y + 1])
		neighbours = neighbours+int(board[x + 1][y + 1])
	}

	return neighbours
}

func (b Board) String() string {
	buffer := bytes.NewBufferString("")
	for _, row := range b.cells {
		buffer.WriteString("\n")
		for _, cell := range row {
			buffer.WriteString(" [")
			buffer.WriteString(strconv.Itoa(int(cell)))
			buffer.WriteString("] ")
		}
	}
	return buffer.String()
}
