package melange

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

func TestHash(t *testing.T) {
	board := NewBoard()
	initialHash := board.CalculateHash()
	assert.Equal(t, initialHash, uint64(7884567450270687084), "Initial hash does not match expected value")

	board.MovePiece(board.NewMove(MoveNormal, E2, E4, Pawn), true)
	board.MovePiece(board.NewMove(MoveNormal, E7, E5, Pawn), false)
	board.MovePiece(board.NewMove(MoveNormal, G1, F3, Knight), true)
	hashAfterMove := board.CalculateHash()
	fmt.Println(board.ToString())
	assert.Equal(t, hashAfterMove, uint64(10074413888190382954), "Hash should change after a move")

	board = NewBoard()
	board.MovePiece(board.NewMove(MoveNormal, G1, F3, Knight), true)
	board.MovePiece(board.NewMove(MoveNormal, E7, E5, Pawn), false)
	board.MovePiece(board.NewMove(MoveNormal, E2, E4, Pawn), true)
	hashAfterDifferentOrder := board.CalculateHash()
	fmt.Println(board.ToString())
	assert.Equal(t, hashAfterMove, hashAfterDifferentOrder, "Hash should be the same for same position reached by different move orders")
}

func BenchmarkHashCalculation(b *testing.B) {
	board := NewBoard()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		board.CalculateHash()
	}
}
