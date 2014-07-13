package conway

import (
	"bytes"
	"log"
	"strconv"
)

type Cell int

const (
	empty = 0
	alive = 1
	dead  = 2

	debug = false
)

type Board struct {
	cells      [][]Cell
	generation int
}

func (this *Board) Init(size int) {
	this.cells = make([][]Cell, size)
	for i := range this.cells {
		this.cells[i] = make([]Cell, size)
	}
	log.Printf("Initialized board with size: [%d][%d]", len(this.cells), len(this.cells[0]))
}

func (this *Board) InitWithSeed(size int, seed [][]Cell) {
	this.Init(size)
	for x, row := range seed {
		for y, cell := range row {
			this.cells[x][y] = cell
		}
	}
	log.Printf("Initialized board with size and seed: [%d][%d] %v", len(this.cells), len(this.cells[0]), seed)
}

func (this *Board) Reaper() {
	debugf("Running reaper. Generation: %d", this.generation)
	this.generation = this.generation+1
	board := copyBoard(this.cells)

	for rowNo, _ := range board {
		for colNo, _ := range board[rowNo] {
			neighbours := countNeighbours(rowNo, colNo, board)
			switch board[rowNo][colNo] {
			case empty:
				if neighbours == 3 {
					debugf("[%d, %d] is alive becase of %d neighbour(s)!", rowNo, colNo, neighbours)
					this.cells[rowNo][colNo] = alive
				} else {
					debugf("[%d, %d] is still dead because of only %d neighbour(s)!", rowNo, colNo, neighbours)
					this.cells[rowNo][colNo] = empty
				}
			case alive:
				if neighbours == 2 || neighbours == 3 {
					debugf("[%d, %d] is still alive because of %d neighbour(s)!", rowNo, colNo, neighbours)
					this.cells[rowNo][colNo] = alive
				} else {
					debugf("[%d, %d] is dying becase of only %d neighbour(s)!", rowNo, colNo, neighbours)
					this.cells[rowNo][colNo] = empty
				}
			}
		}
	}
}

func (this *Board) cell(x int, y int) Cell {
	return this.cells[x][y]
}

func copyBoard(cells [][]Cell) [][]Cell {
	copiedCells := make([][]Cell, len(cells))
	copy(copiedCells, cells)
	for rowNo, row := range cells {
		copiedCells[rowNo] = make([]Cell, len(row))
		copy(copiedCells[rowNo], row)
	}
	return copiedCells
}

func countNeighbours(rowNo int, colNo int, cells [][]Cell) int {
	neighbours := 0

	if (rowNo > 0 && colNo > 0) && (rowNo < len(cells)-1 && colNo < len(cells)-1) {
		//Top row
		neighbours = neighbours+int(cells[rowNo - 1][colNo - 1])
		neighbours = neighbours+int(cells[rowNo - 1][colNo])
		neighbours = neighbours+int(cells[rowNo - 1][colNo + 1])
		//Middle
		neighbours = neighbours+int(cells[rowNo][colNo - 1])
		neighbours = neighbours+int(cells[rowNo][colNo + 1])
		//Bottom row
		neighbours = neighbours+int(cells[rowNo + 1][colNo - 1])
		neighbours = neighbours+int(cells[rowNo + 1][colNo])
		neighbours = neighbours+int(cells[rowNo + 1][colNo + 1])
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

func debugf(format string, v ...interface{}) {
	if debug {
		log.Printf(format, v)
	}
}
