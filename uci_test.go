package muabdib

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestUCIPositionStartpos(t *testing.T) {
	// Reset
	ProcessUciCommand("ucinewgame")
	// Apply start position without moves
	ProcessUciCommand("position startpos")
	b := GetCurrentBoard()

	// Compare with a fresh NewBoard
	ref := NewBoard()
	assert.DeepEqual(t, b.WhitePieces, ref.WhitePieces)
	assert.DeepEqual(t, b.BlackPieces, ref.BlackPieces)
	assert.Equal(t, b.WhiteToMove, true)
	assert.Equal(t, b.Castling, WhiteKingSide|WhiteQueenSide|BlackKingSide|BlackQueenSide)
	assert.Equal(t, b.EnPassant, uint8(0))
	assert.Equal(t, b.HalfMove, uint32(0))
	assert.Equal(t, b.FullMove, uint32(1))
}

func TestUCIPositionStartposWithMoves(t *testing.T) {
	ProcessUciCommand("ucinewgame")
	ProcessUciCommand("position startpos moves e2e4 e7e5 g1f3")
	b := GetCurrentBoard()

	// After e2e4 e7e5 g1f3
	// White knight on f3, white pawn on e4, black pawn on e5
	piece, white := b.PieceAtSquare(F3)
	assert.Equal(t, piece, Knight)
	assert.Equal(t, white, true)

	piece, white = b.PieceAtSquare(E4)
	assert.Equal(t, piece, Pawn)
	assert.Equal(t, white, true)

	piece, white = b.PieceAtSquare(E5)
	assert.Equal(t, piece, Pawn)
	assert.Equal(t, white, false)

	// Side to move should be Black
	assert.Equal(t, b.WhiteToMove, false)
	// Fullmove should be 2 (after Black's move)
	assert.Equal(t, b.FullMove, uint32(2))
}

func TestUCIPositionFenWithMoves(t *testing.T) {
	ProcessUciCommand("ucinewgame")
	// Start from initial FEN using fen syntax
	cmd := "position fen rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1 moves e2e4 c7c5 b1c3"
	ProcessUciCommand(cmd)
	b := GetCurrentBoard()

	// Pieces moved accordingly
	piece, white := b.PieceAtSquare(E4)
	assert.Equal(t, piece, Pawn)
	assert.Equal(t, white, true)

	piece, white = b.PieceAtSquare(C5)
	assert.Equal(t, piece, Pawn)
	assert.Equal(t, white, false)

	piece, white = b.PieceAtSquare(C3)
	assert.Equal(t, piece, Knight)
	assert.Equal(t, white, true)

	// Now Black to move (three moves made)
	assert.Equal(t, b.WhiteToMove, false)
	// Fullmove should be 2 as black has moved once
	assert.Equal(t, b.FullMove, uint32(2))
}
