package muabdib

import (
	"math/bits"
)

var PerftMoveMap = make(map[Move]PerftResult)

// Perft recursively counts the number of leaf nodes reachable within a given depth.
// It generates only legal moves (i.e. moves that do not leave the moving side in check).
// depth == 1 returns the number of legal moves in the current position.
func (b *Board) Perft(depth int, start bool) PerftResult {
	if depth == 0 {
		return PerftResult{Nodes: 1}
	}

	moves := b.GetLegalMoves() // pseudo legal (king safety filtered partly for king moves but not for discovered checks)
	// We must filter moves that leave own king in check.
	res := PerftResult{}
	if depth == 1 {
		for _, m := range moves {
			if b.isMoveLegal(m) {
				res.Nodes++
				if m.IsCapture() {
					res.Captures++
					// Detect "en passant" capture
					fromBB := m.GetFrom64()
					piece, _ := b.PieceAtSquare(fromBB)
					if piece == Pawn {
						toBB := m.GetTo64()
						// Si no hay pieza en el destino pero es captura, es en passant
						if (b.AllPieces() & toBB) == 0 {
							res.EnPassants++
						}
					}
				}
				// Detectar jaques
				state := b.perftMakeMove(m)
				// Tras hacer el movimiento, WhiteToMove indica el lado que debe responder.
				// Si el rey de ese lado está siendo atacado, el movimiento ha dado jaque.
				var kingBB uint64
				kingIsWhite := b.WhiteToMove // color del rey potencialmente en jaque
				if kingIsWhite {
					kingBB = b.WhitePieces.King
				} else {
					kingBB = b.BlackPieces.King
				}
				if kingBB != 0 { // seguridad
					kingSq := uint8(bits.TrailingZeros64(kingBB))
					row := int8(kingSq / 8)
					col := int8(kingSq % 8)
					if b.SquareAttacked(row, col, kingIsWhite) {
						res.Checks++
					}
				}
				b.unmakeMove(state)

				// Detectar promociones
				fromBB := m.GetFrom64()
				piece, _ := b.PieceAtSquare(fromBB)
				if piece == Pawn {
					toBB := m.GetTo64()
					toSq := uint8(bits.TrailingZeros64(toBB))
					toRow := toSq / 8
					if toRow == 0 || toRow == 7 {
						res.Promotions++
					}
				}
				if m.Type == MoveKingCastle || m.Type == MoveQueenCastle {
					res.Castles++
				}
			}
		}
		return res
	}
	for _, m := range moves {
		if !b.isMoveLegal(m) {
			continue
		}
		// Clone here is slover, so we use make/unmake
		state := b.perftMakeMove(m)
		deeperRes := b.Perft(depth-1, false)
		res.Add(deeperRes)
		b.unmakeMove(state)

		if start {
			PerftMoveMap[m] = deeperRes
		}

		// copy := b.Clone()
		// copy.perftMakeMove(m)
		// deeperRes := copy.Perft(depth - 1)
		// res.Nodes += deeperRes.Nodes
		// res.Captures += deeperRes.Captures
	}
	return res
}

type PerftResult struct {
	Nodes      int
	Captures   int
	EnPassants int
	Castles    int
	Promotions int
	Checks     int
}

func (r *PerftResult) Add(other PerftResult) {
	r.Nodes += other.Nodes
	r.Captures += other.Captures
	r.EnPassants += other.EnPassants
	r.Castles += other.Castles
	r.Promotions += other.Promotions
	r.Checks += other.Checks
}

// Perft helper for tests keeping previous API style.
func Perft(board *Board, depth int) PerftResult {
	return board.Perft(depth, true)
}

// moveState stores the information needed to undo a move quickly.
type moveState struct {
	move            Move
	captured        Piece // 0 if none
	capturedIsWhite bool
	whitePieces     Pieces
	blackPieces     Pieces
	castling        CastleRights
	enPassant       uint8
	whiteToMove     bool
}

// perftMakeMove applies a move (already assumed pseudo-legal) and returns the previous state.
func (b *Board) perftMakeMove(m Move) moveState {
	st := moveState{
		move:        m,
		whitePieces: b.WhitePieces,
		blackPieces: b.BlackPieces,
		castling:    b.Castling,
		enPassant:   b.EnPassant,
		whiteToMove: b.WhiteToMove,
	}

	fromBB := m.GetFrom64()
	toBB := m.GetTo64()
	piece, isWhite := b.PieceAtSquare(fromBB)
	// Detect capture (including en passant) before modifying piece sets
	var capturedPiece Piece
	var capturedIsWhite bool
	if m.IsCapture() {
		// Normal capture: piece sitting on destination
		if (b.AllPieces() & toBB) != 0 {
			capturedPiece, capturedIsWhite = b.PieceAtSquare(toBB)
		} else if piece == Pawn { // possible en passant
			if isWhite { // capture black pawn behind
				capturedPiece = Pawn
				capturedIsWhite = false
				toBehind := toBB >> 8
				b.CapturePiece(toBehind, false)
			} else {
				capturedPiece = Pawn
				capturedIsWhite = true
				toBehind := toBB << 8
				b.CapturePiece(toBehind, true)
			}
		}
	}
	if capturedPiece != 0 && (b.AllPieces()&toBB) != 0 { // normal capture remove directly
		b.CapturePiece(toBB, capturedIsWhite)
	}

	// Move the piece
	b.MovePiece(m, isWhite, piece)

	st.captured = capturedPiece
	st.capturedIsWhite = capturedIsWhite
	return st
}

// unmakeMove restores the board to the provided state (inverse of makeMove)
func (b *Board) unmakeMove(st moveState) {
	// Restore bulk state first
	b.WhitePieces = st.whitePieces
	b.BlackPieces = st.blackPieces
	b.Castling = st.castling
	b.EnPassant = st.enPassant
	b.WhiteToMove = st.whiteToMove
}

// isMoveLegal checks if executing m leaves own king in check.
func (b *Board) isMoveLegal(m Move) bool {
	copy := b.Clone()
	st := copy.perftMakeMove(m)
	// Check if previous side's king is attacked.
	var king uint64
	if st.whiteToMove { // white moved
		king = copy.WhitePieces.King
		// if kingBB == 0 { // should not happen
		// 	// b.unmakeMove(st)
		// 	return false
		// }
		kingSq := uint8(bits.TrailingZeros64(king))
		row := int8(kingSq / 8)
		col := int8(kingSq % 8)
		inCheck := copy.SquareAttacked(row, col, true)
		// b.unmakeMove(st)
		return !inCheck
	} else { // black moved
		king = copy.BlackPieces.King
		// if kingBB == 0 { // should not happen
		// 	// b.unmakeMove(st)
		// 	return false
		// }
		kingSq := uint8(bits.TrailingZeros64(king))
		row := int8(kingSq / 8)
		col := int8(kingSq % 8)
		inCheck := copy.SquareAttacked(row, col, false)
		// b.unmakeMove(st)
		return !inCheck
	}
}
