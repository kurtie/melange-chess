package melange

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestEvalMaterial(t *testing.T) {
	board := NewBoard()

	totalMaterial := CpPawn*8 + CpKnight*2 + CpBishop*2 + CpRook*2 + CpQueen + CpKing

	assert.Equal(t, board.WhitePieces.evalMaterial(), totalMaterial)
	assert.Equal(t, board.BlackPieces.evalMaterial(), totalMaterial)
	assert.Equal(t, board.Evaluate(), 0)
}

func TestEvalMaterial2(t *testing.T) {
	board := NewBoard()
	board.SetFen("8/2p5/3p4/KP5r/1R3p1k/8/4P13/8 w - - 0 1 ")

	totalMaterialWhite := CpPawn*2 + CpRook + CpKing
	totalMaterialBlack := CpPawn*3 + CpRook + CpKing

	assert.Equal(t, board.WhitePieces.evalMaterial(), totalMaterialWhite)
	assert.Equal(t, board.BlackPieces.evalMaterial(), totalMaterialBlack)
}

func TestEvalMaterial3(t *testing.T) {
	board := NewBoard()
	board.SetFen("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1")

	totalMaterialWhite := CpPawn*8 + CpKnight*2 + CpBishop*2 + CpRook*2 + CpQueen + CpKing
	totalMaterialBlack := CpPawn*7 + CpKnight*2 + CpBishop*2 + CpRook*2 + CpQueen + CpKing

	assert.Equal(t, board.WhitePieces.evalMaterial(), totalMaterialWhite)
	assert.Equal(t, board.BlackPieces.evalMaterial(), totalMaterialBlack)
}

func TestEvalPosition(t *testing.T) {
	board := NewBoard()
	assert.Equal(t, board.WhitePieces.evalPositions(true), -95)
	assert.Equal(t, board.BlackPieces.evalPositions(false), -95)
}

func TestEvalPosition2(t *testing.T) {
	board := NewBoard()
	board.SetFen("8/2p5/3p4/KP5r/1R3p1k/8/4P13/8 w - - 0 1 ")

	scoreWhite := PosPawn[toIdxSym(B5)] + PosPawn[toIdxSym(E2)] + PosRook[toIdxSym(B4)] + PosKingMiddle[toIdxSym(A5)]
	scoreBlack := PosPawn[toIdx(C7)] + PosPawn[toIdx(D6)] + PosPawn[toIdx(F4)] + PosRook[toIdx(H5)] + PosKingMiddle[toIdx(H4)]

	assert.Equal(t, board.WhitePieces.evalPositions(true), scoreWhite)
	assert.Equal(t, board.BlackPieces.evalPositions(false), scoreBlack)
}

func TestEvalPosition3(t *testing.T) {
	board := NewBoard()
	board.SetFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - ")

	assert.Equal(t, board.WhitePieces.evalPositions(true), 130)
	assert.Equal(t, board.BlackPieces.evalPositions(false), 25)
}

func TestFullEval(t *testing.T) {
	board := NewBoard()
	assert.Equal(t, board.Evaluate(), 0)
}

func TestFullEval2(t *testing.T) {
	board := NewBoard()
	board.SetFen("8/2p5/3p4/KP5r/1R3p1k/8/4P13/8 w - - 0 1 ")
	assert.Equal(t, board.Evaluate(), -130)
}

func TestFullEval3(t *testing.T) {
	board := NewBoard()
	board.SetFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - ")
	assert.Equal(t, board.Evaluate(), 105)
}
