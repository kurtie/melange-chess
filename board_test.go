package muabdib

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestMovePiece_WhitePawn(t *testing.T) {
	board := NewBoard()
	// e2 to e3: e2 = 12, e3 = 20
	move := board.NewMove(MoveNormal, E2, E3, Pawn)
	board.MovePiece(move, true)
	piece, isWhite := board.PieceAtSquare(E3)
	assert.Assert(t, piece == Pawn && isWhite, "Expected white pawn at e3 after move")
	assert.Assert(t, board.WhitePieces.Pawns&E3 != 0, "White pawn should be at e3 after move")
}

func TestMovePiece_BlackPawn(t *testing.T) {
	board := NewBoard()
	// e7 to e6: e7 = 52, e6 = 44
	move := board.NewMove(MoveNormal, E7, E6, Pawn)
	board.MovePiece(move, false)
	piece, isWhite := board.PieceAtSquare(E6)
	assert.Assert(t, piece == Pawn && !isWhite, "Expected black pawn at e6 after move")
	assert.Assert(t, board.BlackPieces.Pawns&E6 != 0, "Black pawn should be at e6 after move")
}

func TestGetLegalMoves_InitialPosition(t *testing.T) {
	board := NewBoard()
	moves := board.GetLegalMoves()
	// In initial position, white has 16 pawn moves (8 single, 8 double) and 4 knight moves
	assert.Assert(t, len(moves) == 20, "Expected 20 legal moves in initial position, got %d", len(moves))
	expected := "a2a3, a2a4, b1a3, b1c3, b2b3, b2b4, c2c3, c2c4, d2d3, d2d4, e2e3, e2e4, f2f3, f2f4, g1f3, g1h3, g2g3, g2g4, h2h3, h2h4"
	assert.Equal(t, moves.ToString(true), expected, "Legal moves do not match expected moves")
}

func checkPosition(moves []*Board, square uint64, expectedPiece Piece, expectedIsWhite bool) bool {
	captureFound := false
	for _, b := range moves {
		piece, isWhite := b.PieceAtSquare(uint64(1) << square) // Check if white pawn is now on b5
		if piece == expectedPiece && expectedIsWhite == isWhite {
			captureFound = true
			break
		}
	}
	return captureFound
}

func TestGetLegalMoves_WhiteCapturePawnRight(t *testing.T) {
	board := &Board{}
	board.WhiteToMove = true
	board.WhitePieces.Pawns = A4
	board.BlackPieces.Pawns = B5
	moves := board.GetLegalMoves()
	// a4 to a5 and a4 captures b5
	expected := "a4a5, a4xb5"
	assert.Equal(t, moves.ToString(true), expected, "Legal moves do not match expected moves")

	// captureFound := checkPosition(moves, 33, Pawn, true)
	// if !captureFound {
	// 	t.Errorf("Expected capture move to b5 not found")
	// }
}

func TestGetLegalMoves_WhiteCapturePawnLeft(t *testing.T) {
	board := &Board{}
	board.WhiteToMove = true
	board.WhitePieces.Pawns = B4
	board.BlackPieces.Pawns = A5
	moves := board.GetLegalMoves()
	// b4 to b5 and b4 captures a5
	expected := "b4b5, b4xa5"
	assert.Equal(t, moves.ToString(true), expected, "Legal moves do not match expected moves")

	// captureFound := checkPosition(moves, 32, Pawn, true)
	// if !captureFound {
	// 	t.Errorf("Expected capture move to a5 not found")
	// }
}

func TestGetLegalMoves_BlackCapturePawnRight(t *testing.T) {
	board := &Board{}
	board.WhiteToMove = false
	board.BlackPieces.Pawns = A5
	board.WhitePieces.Pawns = B4
	moves := board.GetLegalMoves()
	// a5 to a4 and a5 captures b4
	expected := "a5a4, a5xb4"
	assert.Equal(t, moves.ToString(true), expected, "Legal moves do not match expected moves")
	// captureFound := checkPosition(moves, 25, Pawn, false)
	// if !captureFound {
	// 	t.Errorf("Expected capture move to b4 not found")
	// }
}

func TestGetLegalMoves_BlackCapturePawnLeft(t *testing.T) {
	board := &Board{}
	board.WhiteToMove = false
	board.BlackPieces.Pawns = B5
	board.WhitePieces.Pawns = A4
	moves := board.GetLegalMoves()
	// b5 to b4 and b5 captures a4
	expected := "b5b4, b5xa4"
	assert.Equal(t, moves.ToString(true), expected, "Legal moves do not match expected moves")
	// captureFound := checkPosition(moves, 24, Pawn, false)
	// if !captureFound {
	// 	t.Errorf("Expected capture move to a4 not found")
	// }
}

func TestGetLegalMoves_KnightCornerSW(t *testing.T) {
	board := &Board{}
	board.WhiteToMove = true
	board.WhitePieces.Knights = B2
	moves := board.GetLegalMoves()
	// Knight can move to a4 (16) and c4 (64), d3 (32), d1 (8)
	expected := "b2a4, b2c4, b2d1, b2d3"
	assert.Equal(t, moves.ToString(true), expected, "Legal moves do not match expected moves")
}

func TestGetLegalMoves_KnightCornerSE(t *testing.T) {
	board := &Board{}
	board.WhiteToMove = false
	board.BlackPieces.Knights = H2
	moves := board.GetLegalMoves()
	// Knight can move to g4 (64), f3 (32), f1 (8)
	expected := "h2f1, h2f3, h2g4"
	assert.Equal(t, moves.ToString(true), expected, "Legal moves do not match expected moves")
}

func TestGetLegalMoves_KnightCornerNW(t *testing.T) {
	board := &Board{}
	board.WhiteToMove = false
	board.BlackPieces.Knights = A8
	moves := board.GetLegalMoves()
	// Knight can move to b6 (16), c7 (32)
	expected := "a8b6, a8c7"
	assert.Equal(t, moves.ToString(true), expected, "Legal moves do not match expected moves")
}

func TestGetLegalMoves_KnightCornerNE(t *testing.T) {
	board := &Board{}
	board.WhiteToMove = true
	board.WhitePieces.Knights = H8
	moves := board.GetLegalMoves()

	// Knight can move to g6 (64), f7 (32)
	expected := "h8f7, h8g6"
	assert.Equal(t, moves.ToString(true), expected, "Legal moves do not match expected moves")
}

func TestGetLegalMoves_Bishop(t *testing.T) {
	board := &Board{}
	board.WhiteToMove = true
	board.WhitePieces.Bishops = E5
	board.BlackPieces.Pawns = G3
	board.WhitePieces.Pawns = D6
	moves := board.GetLegalMoves()
	// Bishop can move to a1, b2, c3, d4, f4, f6, g7, h8, and pawn can move d6d7
	expected := "d6d7, e5a1, e5b2, e5c3, e5d4, e5f4, e5f6, e5g7, e5h8, e5xg3"
	assert.Equal(t, moves.ToString(true), expected, "Legal moves do not match expected moves")
}

func TestGetLegalMoves_Rook(t *testing.T) {
	board := &Board{}
	board.WhiteToMove = false
	board.BlackPieces.Rooks = E5
	board.BlackPieces.Pawns = E3
	board.WhitePieces.Pawns = B5
	moves := board.GetLegalMoves()
	// Rook can move to c5, d5, e4, e6, e7, e8, f5, g5, h5, b5 (capture), and pawn can move e3e2
	expected := "e3e2, e5c5, e5d5, e5e4, e5e6, e5e7, e5e8, e5f5, e5g5, e5h5, e5xb5"
	assert.Equal(t, moves.ToString(true), expected, "Legal moves do not match expected moves")
}

func TestGetLegalMoves_Queen(t *testing.T) {
	board := &Board{}
	board.WhiteToMove = true
	board.WhitePieces.Queens = E5
	board.BlackPieces.Pawns = E3
	board.WhitePieces.Pawns = B5
	moves := board.GetLegalMoves()
	// Queen can move to 22 places and 1 capture,  and pawn can move b5b6
	expected := "b5b6, e5a1, e5b2, e5b8, e5c3, e5c5, e5c7, e5d4, e5d5, e5d6, e5e4, e5e6, e5e7, e5e8, e5f4, e5f5, e5f6, e5g3, e5g5, e5g7, e5h2, e5h5, e5h8, e5xe3"
	assert.Equal(t, moves.ToString(true), expected, "Legal moves do not match expected moves")
}

func TestGetLegalMoves_King(t *testing.T) {
	board := &Board{}
	board.WhiteToMove = true
	board.WhitePieces.King = E5
	board.BlackPieces.Pawns = F5
	board.WhitePieces.Pawns = D5
	moves := board.GetLegalMoves()
	// King can move to d4, d6, e6, f4, f6, capture f5 and pawn can move d5d6
	expected := "d5d6, e5d4, e5d6, e5e6, e5f4, e5f6, e5xf5"
	assert.Equal(t, moves.ToString(true), expected, "Legal moves do not match expected moves")
}

func TestSquareAttackedByKnights(t *testing.T) {
	board := &Board{}
	board = &Board{}
	board.WhiteToMove = true
	board.WhitePieces.King = E5
	board.BlackPieces.Knights = F5
	board.WhitePieces.Pawns = D5
	moves := board.GetLegalMoves()
	// Pawn can move 1, King can move 5 (2 attacked by knight)
	expected := "d5d6, e5e4, e5e6, e5f4, e5f6, e5xf5"
	assert.Equal(t, moves.ToString(true), expected, "Legal moves do not match expected moves")
}

func TestSquareAttackedByBishops(t *testing.T) {
	board := &Board{}
	board = &Board{}
	board.WhiteToMove = true
	board.WhitePieces.King = E5
	board.BlackPieces.Bishops = F5
	board.WhitePieces.Pawns = D5
	legalMoves := board.GetLegalMoves()
	// Pawn can move 1, King can move 5 (2 attacked by bishop)
	expected := "d5d6, e5d4, e5d6, e5f4, e5f6, e5xf5"
	assert.Equal(t, legalMoves.ToString(true), expected, "Legal moves do not match expected moves")
}

func TestSquareAttackedByRooks(t *testing.T) {
	board := &Board{}
	board = &Board{}
	board.WhiteToMove = true
	board.WhitePieces.King = E5
	board.BlackPieces.Rooks = F5
	board.WhitePieces.Pawns = D5
	legalMoves := board.GetLegalMoves()
	// Pawn can move 1, King can move 5 (2 attacked by rook)
	expected := "d5d6, e5d4, e5d6, e5e4, e5e6, e5xf5"
	assert.Equal(t, legalMoves.ToString(true), expected, "Legal moves do not match expected moves")
}

func TestSquareAttackedByQueens(t *testing.T) {
	board := &Board{}
	board = &Board{}
	board.WhiteToMove = true
	board.WhitePieces.King = E5
	board.BlackPieces.Queens = F5
	board.WhitePieces.Pawns = D5
	legalMoves := board.GetLegalMoves()
	// Pawn can move 1, King can move 3 (4 attacked by queen)
	expected := "d5d6, e5d4, e5d6, e5xf5"
	assert.Equal(t, legalMoves.ToString(true), expected, "Legal moves do not match expected moves")
}

func TestSquareAttackedByKing(t *testing.T) {
	board := &Board{}
	board = &Board{}
	board.WhiteToMove = true
	board.WhitePieces.King = E5
	board.BlackPieces.King = G5
	board.WhitePieces.Pawns = D5
	legalMoves := board.GetLegalMoves()
	// Pawn can move 1, King can move 4
	expected := "d5d6, e5d4, e5d6, e5e4, e5e6"
	assert.Equal(t, legalMoves.ToString(true), expected, "Legal moves do not match expected moves")
}

func TestWhiteKingSideCastle(t *testing.T) {
	board := &Board{}
	board.WhiteToMove = true
	board.WhitePieces.King = E1
	board.WhitePieces.Rooks = H1
	board.Castling = WhiteKingSide
	legalMoves := board.GetLegalMoves()
	// King can castle king-side
	expected := "O-O, e1d1, e1d2, e1e2, e1f1, e1f2, h1f1, h1g1, h1h2, h1h3, h1h4, h1h5, h1h6, h1h7, h1h8"
	assert.Equal(t, legalMoves.ToString(true), expected, "Legal moves do not match expected moves")
}

func TestWhiteQueenSideCastle(t *testing.T) {
	board := &Board{}
	board.WhiteToMove = true
	board.WhitePieces.King = E1
	board.WhitePieces.Rooks = A1
	board.Castling = WhiteQueenSide
	legalMoves := board.GetLegalMoves()
	// King can castle queen-side
	expected := "O-O-O, a1a2, a1a3, a1a4, a1a5, a1a6, a1a7, a1a8, a1b1, a1c1, a1d1, e1d1, e1d2, e1e2, e1f1, e1f2"
	assert.Equal(t, legalMoves.ToString(true), expected, "Legal moves do not match expected moves")
}

func TestBlackKingSideCastle(t *testing.T) {
	board := &Board{}
	board.WhiteToMove = false
	board.BlackPieces.King = E8
	board.BlackPieces.Rooks = H8
	board.Castling = BlackKingSide
	legalMoves := board.GetLegalMoves()
	// King can castle king-side
	expected := "O-O, e8d7, e8d8, e8e7, e8f7, e8f8, h8f8, h8g8, h8h1, h8h2, h8h3, h8h4, h8h5, h8h6, h8h7"
	assert.Equal(t, legalMoves.ToString(true), expected, "Legal moves do not match expected moves")
}

func TestBlackQueenSideCastle(t *testing.T) {
	board := &Board{}
	board.WhiteToMove = false
	board.BlackPieces.King = E8
	board.BlackPieces.Rooks = A8
	board.Castling = BlackQueenSide
	legalMoves := board.GetLegalMoves()
	// King can castle queen-side
	expected := "O-O-O, a8a1, a8a2, a8a3, a8a4, a8a5, a8a6, a8a7, a8b8, a8c8, a8d8, e8d7, e8d8, e8e7, e8f7, e8f8"
	assert.Equal(t, legalMoves.ToString(true), expected, "Legal moves do not match expected moves")
}

func TestEnPassantGenerationWhite(t *testing.T) {
	// Posición: peón blanco en e5 (e5 index 28), peón negro acaba de mover d7-d5 dejando d6 como en passant (d6 index 35)
	board := &Board{}
	board.WhiteToMove = true
	board.WhitePieces.Pawns = E5
	board.BlackPieces.Pawns = D5 // peón negro en d5 (index 35-? d5 is 35)
	// Simular que el último movimiento fue d7-d5 => target en passant es d6 (index 43? recalculamos)
	// Indices: a1=0 => d5 = (fila 5-1=4)*8 + (col d=3) = 4*8+3=35 correcto. d6 = (fila 6-1=5)*8+3=43
	board.EnPassant = 43
	moves := board.GetLegalMoves()
	expected := "e5e6, e5xd6"
	assert.Equal(t, moves.ToString(true), expected, "Legal moves do not match expected moves")
}

func TestEnPassantExecutionWhite(t *testing.T) {
	// Configurar posición similar y ejecutar en passant
	board := &Board{}
	board.WhiteToMove = true
	board.WhitePieces.Pawns = E5
	board.BlackPieces.Pawns = D5
	board.EnPassant = 43 // d6
	// Crear movimiento e5xd6 (en passant)
	move := board.NewMove(MoveCapture, E5, D6, Pawn)
	// Marcar que el tablero tenía EnPassant previo
	move.EnPassant = 0 // Después del movimiento se limpia
	board.MovePiece(move, true)
	// Peón negro en d5 debería haber sido capturado
	assert.Assert(t, board.BlackPieces.Pawns&D5 == 0, "Black pawn on d5 should be captured via en passant")
}

func TestEnPassantGenerationBlack(t *testing.T) {
	// Peón negro en d4 (d4 index 27), peón blanco acaba de mover e2-e4 => target e3 (index 20)
	board := &Board{}
	board.WhiteToMove = false
	board.BlackPieces.Pawns = D4
	board.WhitePieces.Pawns = E4
	// e4 index: (fila 4-1=3)*8 + 4? file e=4 => 3*8+4=28 (coincide con E5 antes) Wait: E4 constant is 28 yes.
	// e3 index: (fila 3-1=2)*8 +4 = 20
	board.EnPassant = 20
	moves := board.GetLegalMoves()
	expected := "d4d3, d4xe3"
	assert.Equal(t, moves.ToString(true), expected, "Legal moves do not match expected moves")
}

func TestEnPassantExecutionBlack(t *testing.T) {
	board := &Board{}
	board.WhiteToMove = false
	board.BlackPieces.Pawns = D4
	board.WhitePieces.Pawns = E4
	board.EnPassant = 20 // e3
	move := board.NewMove(MoveCapture, D4, E3, Pawn)
	board.MovePiece(move, false)
	assert.Assert(t, board.WhitePieces.Pawns&E4 == 0, "White pawn on e4 should be captured via en passant")
}

func TestWhitePromotionGeneration(t *testing.T) {
	board := &Board{}
	board.WhiteToMove = true
	// White pawn on a7 (ready to promote), and black piece on b8 to allow capture promotion
	board.WhitePieces.Pawns = A7
	board.BlackPieces.Knights = B8
	moves := board.GetLegalMoves()
	// From a7 -> a8 promotions (4) and a7xb8 promotions (4); sort expected lexicographically
	expectedSlice := []string{"a7a8=B", "a7a8=N", "a7a8=Q", "a7a8=R", "a7xb8=B", "a7xb8=N", "a7xb8=Q", "a7xb8=R"}
	got := moves.ToStringArraySorted()
	assert.DeepEqual(t, got, expectedSlice)
}

func TestBlackPromotionGeneration(t *testing.T) {
	board := &Board{}
	board.WhiteToMove = false
	// Black pawn on h2 (index 15) ready to promote moving to h1, white piece on g1 for capture promotions
	board.BlackPieces.Pawns = H2
	board.WhitePieces.Knights = G1
	moves := board.GetLegalMoves()
	expectedSlice := []string{"h2h1=B", "h2h1=N", "h2h1=Q", "h2h1=R", "h2xg1=B", "h2xg1=N", "h2xg1=Q", "h2xg1=R"}
	got := moves.ToStringArraySorted()
	assert.DeepEqual(t, got, expectedSlice)
}
