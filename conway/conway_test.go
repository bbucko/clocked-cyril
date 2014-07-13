package conway

import (
	"fmt"
	"testing"
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

func TestReaperWithBlock(t *testing.T) {
	boardSize := 10

	board := Board{}
	seed := createSeed(boardSize)

	seed[4][4] = alive
	seed[4][5] = alive
	seed[5][5] = alive
	seed[5][4] = alive

	board.InitWithSeed(boardSize, seed)
	board.Reaper()
	if !(board.cell(4, 4) == alive && board.cell(5, 5) == alive && board.cell(4, 5) == alive && board.cell(5, 4) == alive) {
		t.Errorf("Board was incorrectly modified by reaper")
	}
}

func TestReaperWithBlinker(t *testing.T) {
	boardSize := 10

	board := Board{}
	seed := createSeed(boardSize)

	seed[4][5] = alive
	seed[5][5] = alive
	seed[6][5] = alive

	board.InitWithSeed(boardSize, seed)
	if !(board.cell(4, 5) == alive && board.cell(5, 5) == alive && board.cell(6, 5) == alive) {
		t.Errorf("Board was incorrectly initialized [%d] [%d] [%d]", board.cell(4, 5), board.cell(5, 5), board.cell(6, 5))
	}

	board.Reaper()
	if !(board.cell(5, 4) == alive && board.cell(5, 5) == alive && board.cell(5, 6) == alive) {
		t.Errorf("Board was incorrectly modified by reaper [%d] [%d] [%d]", board.cell(5, 4), board.cell(5, 5), board.cell(5, 6))
	}
}

func BenchmarkBlinker(b *testing.B) {
	boardSize := 10

	board := Board{}
	seed := createSeed(boardSize)

	seed[4][5] = alive
	seed[5][5] = alive
	seed[6][5] = alive

	board.InitWithSeed(boardSize, seed)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		board.Reaper()
	}
}

func BenchmarkBlock(b *testing.B) {
	boardSize := 10

	board := Board{}
	seed := createSeed(boardSize)

	seed[4][4] = alive
	seed[4][5] = alive
	seed[5][5] = alive
	seed[5][4] = alive

	board.InitWithSeed(boardSize, seed)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		board.Reaper()
	}
}

func BenchmarkFullBoard(b *testing.B) {
	boardSize := 10

	board := Board{}
	seed := createSeed(boardSize)

	for x, _ := range seed {
		for y, _ := range seed[x] {
			seed[x][y] = alive
		}
	}

	board.InitWithSeed(boardSize, seed)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		board.Reaper()
	}
}

func createSeed(size int) [][]Cell {
	seed := make([][]Cell, size)
	for i := range seed {
		seed[i] = make([]Cell, size)
	}
	return seed
}
