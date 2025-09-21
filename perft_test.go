package muabdib

import (
	"fmt"
	"sort"
	"testing"

	"gotest.tools/v3/assert"
)

func TestPerft(t *testing.T) {
	board := NewBoard()
	// Initial position
	res := Perft(board, 1)
	assert.Equal(t, res.Nodes, 20)
	assert.Equal(t, res.Captures, 0)
	assert.Equal(t, res.EnPassants, 0)
	assert.Equal(t, res.Checks, 0)
	assert.Equal(t, res.Promotions, 0)
	assert.Equal(t, res.Castles, 0)

	res = Perft(board, 2)
	assert.Equal(t, res.Nodes, 400)
	assert.Equal(t, res.Captures, 0)
	assert.Equal(t, res.EnPassants, 0)
	assert.Equal(t, res.Checks, 0)
	assert.Equal(t, res.Promotions, 0)
	assert.Equal(t, res.Castles, 0)

	res = Perft(board, 3)
	assert.Equal(t, res.Nodes, 8902)
	assert.Equal(t, res.Captures, 34)
	assert.Equal(t, res.EnPassants, 0)
	assert.Equal(t, res.Checks, 12)
	assert.Equal(t, res.Promotions, 0)
	assert.Equal(t, res.Castles, 0)

	res = Perft(board, 4)
	assert.Equal(t, res.Nodes, 197281)
	assert.Equal(t, res.Captures, 1576)
	assert.Equal(t, res.EnPassants, 0)
	assert.Equal(t, res.Checks, 469)
	assert.Equal(t, res.Promotions, 0)
	assert.Equal(t, res.Castles, 0)

	res = Perft(board, 5)
	assert.Equal(t, res.Nodes, 4865609)
	assert.Equal(t, res.Captures, 82719)
	assert.Equal(t, res.EnPassants, 258)
	assert.Equal(t, res.Checks, 27351)
	assert.Equal(t, res.Promotions, 0)
	assert.Equal(t, res.Castles, 0)

	res = Perft(board, 6)
	assert.Equal(t, res.Nodes, 119060324)
	assert.Equal(t, res.Captures, 2812008)
	assert.Equal(t, res.EnPassants, 5248)
	assert.Equal(t, res.Checks, 809099)
	assert.Equal(t, res.Promotions, 0)
	assert.Equal(t, res.Castles, 0)
}

func TestPerftPosition2(t *testing.T) {
	board := NewBoard()
	board.SetFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -")
	res := Perft(board, 1)
	assert.Equal(t, res.Nodes, 48)
	assert.Equal(t, res.Captures, 8)
	assert.Equal(t, res.EnPassants, 0)
	assert.Equal(t, res.Checks, 0)
	assert.Equal(t, res.Promotions, 0)
	assert.Equal(t, res.Castles, 2)

	res = Perft(board, 2)
	assert.Equal(t, res.Nodes, 2039)
	assert.Equal(t, res.Captures, 351)
	assert.Equal(t, res.EnPassants, 1)
	assert.Equal(t, res.Checks, 3)
	assert.Equal(t, res.Promotions, 0)
	assert.Equal(t, res.Castles, 91)

	res = Perft(board, 3)
	assert.Equal(t, res.Nodes, 97862)
	assert.Equal(t, res.Captures, 17102)
	assert.Equal(t, res.EnPassants, 45)
	assert.Equal(t, res.Checks, 993)
	assert.Equal(t, res.Promotions, 0)
	assert.Equal(t, res.Castles, 3162)

	res = Perft(board, 4)
	assert.Equal(t, res.Nodes, 4085603)
	assert.Equal(t, res.Captures, 757163)
	assert.Equal(t, res.EnPassants, 1929)
	assert.Equal(t, res.Checks, 25523)
	assert.Equal(t, res.Promotions, 15172)
	assert.Equal(t, res.Castles, 128013)

	res = Perft(board, 5)
	assert.Equal(t, res.Nodes, 193690690)
	assert.Equal(t, res.Captures, 35043416)
	assert.Equal(t, res.EnPassants, 73365)
	assert.Equal(t, res.Checks, 3309887)
	assert.Equal(t, res.Promotions, 8392)
	assert.Equal(t, res.Castles, 4993637)
}

// Imprimir el mapa ordenado por k.ToSimpleString()
func printMoveResults() {
	moveKeys := make([]Move, 0, len(PerftMoveMap))
	for k := range PerftMoveMap {
		moveKeys = append(moveKeys, k)
	}
	sort.Slice(moveKeys, func(i, j int) bool {
		return moveKeys[i].ToSimpleString() < moveKeys[j].ToSimpleString()
	})
	for _, k := range moveKeys {
		fmt.Println(k.ToSimpleString(), PerftMoveMap[k].Nodes)
	}
}

func TestPerftPosition3(t *testing.T) {
	board := NewBoard()
	board.SetFen("8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1")
	res := Perft(board, 1)
	assert.Equal(t, res.Nodes, 14)
	assert.Equal(t, res.Captures, 1)
	assert.Equal(t, res.EnPassants, 0)
	assert.Equal(t, res.Checks, 2)
	assert.Equal(t, res.Promotions, 0)

	res = Perft(board, 2)
	assert.Equal(t, res.Nodes, 191)
	assert.Equal(t, res.Captures, 14)
	assert.Equal(t, res.EnPassants, 0)
	assert.Equal(t, res.Checks, 10)
	assert.Equal(t, res.Promotions, 0)

	res = Perft(board, 3)
	assert.Equal(t, res.Nodes, 2812)
	assert.Equal(t, res.Captures, 209)
	assert.Equal(t, res.EnPassants, 2)
	assert.Equal(t, res.Checks, 267)
	assert.Equal(t, res.Promotions, 0)

	res = Perft(board, 4)
	assert.Equal(t, res.Nodes, 43238)
	assert.Equal(t, res.Captures, 3348)
	assert.Equal(t, res.EnPassants, 123)
	assert.Equal(t, res.Checks, 1680)
	assert.Equal(t, res.Promotions, 0)

	res = Perft(board, 5)
	assert.Equal(t, res.Nodes, 674624)
	assert.Equal(t, res.Captures, 52051)
	assert.Equal(t, res.EnPassants, 1165)
	assert.Equal(t, res.Checks, 52950)
	assert.Equal(t, res.Promotions, 0)

	res = Perft(board, 6)
	assert.Equal(t, res.Nodes, 11030083)
	assert.Equal(t, res.Captures, 940350)
	assert.Equal(t, res.EnPassants, 33325)
	assert.Equal(t, res.Checks, 452473)
	assert.Equal(t, res.Promotions, 7552)

	res = Perft(board, 7)
	assert.Equal(t, res.Nodes, 178633661)
	assert.Equal(t, res.Captures, 14519036)
	assert.Equal(t, res.EnPassants, 294874)
	assert.Equal(t, res.Checks, 12797406)
	assert.Equal(t, res.Promotions, 140024)
}
