package conway

import (
	"testing"
	"fmt"
)

func TestToString(t *testing.T) {
	boardSize := 10

	board := Board{}
	if fmt.Sprintf("%v", board) != "" {
		t.Errorf("String should return %s but returns %s", "[]", fmt.Sprintf("%v", board))
	}

	board.Init(boardSize)

	if fmt.Sprintf("%v", board) == "" {
		t.Errorf("String should return %s but returns %s", "[]", fmt.Sprintf("%v", board))
	}
}

func TestBasicSetup(t *testing.T) {
	boardSize := 10

	board := Board{}
	board.Init(boardSize)
	if len(board.cells) != boardSize {
		t.Errorf("Board size should be %d and is %d", boardSize, len(board.cells))
	}
}

func TestBasicSetupWithSeed(t *testing.T) {
	boardSize := 10

	board := Board{}
	seed := createSeed(boardSize)

	seed[4][4] = alive
	seed[4][5] = alive
	seed[5][5] = alive
	seed[5][4] = alive

	board.InitWithSeed(boardSize, seed)
	if len(board.cells) != boardSize {
		t.Errorf("Board size should be %d and is %d", boardSize, len(board.cells))
	}
	if board.cell(4, 4) != alive || board.cell(5, 5) != alive {
		t.Errorf("Board was incorrectly created from seed")
	}
}

func TestReaperWithCellsInTheCenter(t *testing.T) {
	boardSize := 10

	board := Board{}
	seed := createSeed(boardSize)

	seed[4][4] = alive
	seed[4][5] = alive
	seed[5][5] = alive
	seed[5][4] = alive

	board.InitWithSeed(boardSize, seed)
	board.Reaper()
	if board.cell(4, 4) != alive || board.cell(5, 5) != alive || board.cell(4, 5) != alive || board.cell(5, 4) != alive {
		t.Errorf("Board was incorrectly modified by reaper")
	}
}

func createSeed(size int) [][]Cell {
	seed := make([][]Cell, size)
	for i := range seed {
		seed[i] = make([]Cell, size)
	}
	return seed
}
