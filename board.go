package muabdib

import (
	"fmt"
	"math"
	"math/bits"
)

type Board struct {
	WhitePieces Pieces
	BlackPieces Pieces
	WhiteToMove bool
	Castling    CastleRights // Bitmask for castling rights
	EnPassant   uint8        // Square index for en passant target (0-63), 0 if none
	HalfMove    uint32       // Halfmove clock for fifty-move rule
	FullMove    uint32       // Fullmove number starting at 1 and incremented after Black's move
}

func NewBoard() *Board {
	return &Board{
		WhitePieces: Pieces{
			Pawns:   A2 | B2 | C2 | D2 | E2 | F2 | G2 | H2,
			Knights: B1 | G1,
			Bishops: C1 | F1,
			Rooks:   A1 | H1,
			Queens:  D1,
			King:    E1,
		},
		BlackPieces: Pieces{
			Pawns:   A7 | B7 | C7 | D7 | E7 | F7 | G7 | H7,
			Knights: B8 | G8,
			Bishops: C8 | F8,
			Rooks:   A8 | H8,
			Queens:  D8,
			King:    E8,
		},
		WhiteToMove: true,
		Castling:    WhiteKingSide | WhiteQueenSide | BlackKingSide | BlackQueenSide,
		EnPassant:   0,
		HalfMove:    0,
		FullMove:    1,
	}
}

func (b *Board) Clone() *Board {
	return &Board{
		WhitePieces: b.WhitePieces,
		BlackPieces: b.BlackPieces,
		WhiteToMove: b.WhiteToMove,
		Castling:    b.Castling,
		EnPassant:   b.EnPassant,
		HalfMove:    b.HalfMove,
		FullMove:    b.FullMove,
	}
}

func (b *Board) NewMove(t MoveType, from uint64, to uint64) Move {
	return Move{
		From:     uint8(bits.TrailingZeros64(from)),
		To:       uint8(bits.TrailingZeros64(to)),
		Type:     t,
		Castling: b.Castling,
		// EnPassant siempre inicia en 0; solo los dobles avances de peón lo establecerán
		EnPassant: 0,
	}
}

func (b *Board) AllPieces() uint64 {
	return b.WhitePieces.Pawns | b.WhitePieces.Knights | b.WhitePieces.Bishops |
		b.WhitePieces.Rooks | b.WhitePieces.Queens | b.WhitePieces.King |
		b.BlackPieces.Pawns | b.BlackPieces.Knights | b.BlackPieces.Bishops |
		b.BlackPieces.Rooks | b.BlackPieces.Queens | b.BlackPieces.King
}

func (b *Board) OccupiedSquares() uint64 {
	return b.AllPieces()
}

func (b *Board) EmptySquares() uint64 {
	return ^b.OccupiedSquares()
}

func (b *Board) WhiteOccupiedSquares() uint64 {
	return b.WhitePieces.Pawns | b.WhitePieces.Knights | b.WhitePieces.Bishops |
		b.WhitePieces.Rooks | b.WhitePieces.Queens | b.WhitePieces.King
}

func (b *Board) BlackOccupiedSquares() uint64 {
	return b.BlackPieces.Pawns | b.BlackPieces.Knights | b.BlackPieces.Bishops |
		b.BlackPieces.Rooks | b.BlackPieces.Queens | b.BlackPieces.King
}

func (b *Board) IsSquareOccupied(square uint64) bool {
	return (b.OccupiedSquares() & square) != 0
}

func (b *Board) IsSquareOccupiedByWhite(square uint64) bool {
	return (b.WhiteOccupiedSquares() & square) != 0
}

func (b *Board) IsSquareOccupiedByBlack(square uint64) bool {
	return (b.BlackOccupiedSquares() & square) != 0
}

func (b *Board) GetPiecesToMove() *Pieces {
	if b.WhiteToMove {
		return &b.WhitePieces
	}
	return &b.BlackPieces
}

func (b *Board) MovePiece(move Move, isWhite bool, pieceType Piece) {
	var pieces *Pieces
	if isWhite {
		pieces = &b.WhitePieces
	} else {
		pieces = &b.BlackPieces
	}

	// Detectar captura en passant antes de mover la pieza: si es captura de peón a casilla EnPassant previa y
	// la casilla destino está vacía
	if pieceType == Pawn && move.Type == MoveCapture {
		destBit := move.GetTo64()
		// La casilla EnPassant válida es la almacenada actualmente en el tablero (del movimiento previo del rival)
		if b.EnPassant == move.To && !b.IsSquareOccupied(destBit) { // casilla destino vacía => en passant
			if isWhite {
				// Captura peón negro que está justo detrás (una fila abajo en términos de bitboard: destino >> 8)
				captured := destBit >> 8
				b.CapturePiece(captured, false)
			} else {
				// Captura peón blanco que está justo detrás (una fila arriba: destino << 8)
				captured := destBit << 8
				b.CapturePiece(captured, true)
			}
		}
	}

	switch pieceType {
	case Pawn:
		pieces.Pawns &= ^move.GetFrom64()
		pieces.Pawns |= move.GetTo64()
	case Knight:
		pieces.Knights &= ^move.GetFrom64()
		pieces.Knights |= move.GetTo64()
	case Bishop:
		pieces.Bishops &= ^move.GetFrom64()
		pieces.Bishops |= move.GetTo64()
	case Rook:
		pieces.Rooks &= ^move.GetFrom64()
		pieces.Rooks |= move.GetTo64()
	case Queen:
		pieces.Queens &= ^move.GetFrom64()
		pieces.Queens |= move.GetTo64()
	case King:
		pieces.King &= ^move.GetFrom64()
		pieces.King |= move.GetTo64()
	}

	b.Castling = move.Castling
	// Si el movimiento es un doble avance de peón establecer EnPassant, si no limpiarlo
	if pieceType == Pawn {
		// Un doble avance se reconoce porque la diferencia de índices es 16 (dos filas) y no es captura
		if (move.Type&MoveCapture) == 0 && (move.From+16 == move.To || move.From == move.To+16) {
			// Casilla intermedia = (from+to)/2
			mid := (uint16(move.From) + uint16(move.To)) / 2
			b.EnPassant = uint8(mid)
		} else {
			b.EnPassant = 0
		}
	} else {
		b.EnPassant = 0
	}
	b.WhiteToMove = !b.WhiteToMove
}

func (b *Board) CapturePiece(square uint64, isWhite bool) {
	if isWhite {
		// Captura pieza blanca
		b.WhitePieces.Pawns &= ^square
		b.WhitePieces.Knights &= ^square
		b.WhitePieces.Bishops &= ^square
		b.WhitePieces.Rooks &= ^square
		b.WhitePieces.Queens &= ^square
	} else {
		// Captura pieza negra
		b.BlackPieces.Pawns &= ^square
		b.BlackPieces.Knights &= ^square
		b.BlackPieces.Bishops &= ^square
		b.BlackPieces.Rooks &= ^square
		b.BlackPieces.Queens &= ^square
	}
}

func (b *Board) PieceAtSquare(square uint64) (Piece, bool) {
	if b.WhitePieces.Pawns&square != 0 {
		return Pawn, true
	} else if b.WhitePieces.Knights&square != 0 {
		return Knight, true
	} else if b.WhitePieces.Bishops&square != 0 {
		return Bishop, true
	} else if b.WhitePieces.Rooks&square != 0 {
		return Rook, true
	} else if b.WhitePieces.Queens&square != 0 {
		return Queen, true
	} else if b.WhitePieces.King&square != 0 {
		return King, true
	} else if b.BlackPieces.Pawns&square != 0 {
		return Pawn, false
	} else if b.BlackPieces.Knights&square != 0 {
		return Knight, false
	} else if b.BlackPieces.Bishops&square != 0 {
		return Bishop, false
	} else if b.BlackPieces.Rooks&square != 0 {
		return Rook, false
	} else if b.BlackPieces.Queens&square != 0 {
		return Queen, false
	} else if b.BlackPieces.King&square != 0 {
		return King, false
	}
	return 0, false
}

func (b *Board) ToString() string {
	board := [8][8]string{}

	for i := 0; i < 64; i++ {
		square := uint64(1) << i
		row := 7 - (i / 8)
		col := i % 8

		res := "."
		if b.WhitePieces.Pawns&square != 0 {
			res = "P"
		} else if b.WhitePieces.Knights&square != 0 {
			res = "N"
		} else if b.WhitePieces.Bishops&square != 0 {
			res = "B"
		} else if b.WhitePieces.Rooks&square != 0 {
			res = "R"
		} else if b.WhitePieces.Queens&square != 0 {
			res = "Q"
		} else if b.WhitePieces.King&square != 0 {
			res = "K"
		} else if b.BlackPieces.Pawns&square != 0 {
			res = "p"
		} else if b.BlackPieces.Knights&square != 0 {
			res = "n"
		} else if b.BlackPieces.Bishops&square != 0 {
			res = "b"
		} else if b.BlackPieces.Rooks&square != 0 {
			res = "r"
		} else if b.BlackPieces.Queens&square != 0 {
			res = "q"
		} else if b.BlackPieces.King&square != 0 {
			res = "k"
		}
		board[row][col] = res
	}

	printPieces := func(bitBoard uint64) string {
		if bitBoard == 0 {
			return "                "
		}
		return fmt.Sprintf("%016x", bitBoard)
	}

	result := ""
	for r, row := range board {
		for _, cell := range row {
			result += cell + " "
		}
		switch r {
		case 0:
			result += "    WP:" + printPieces(b.WhitePieces.Pawns) + "    " + "BP:" + printPieces(b.BlackPieces.Pawns)
		case 1:
			result += "    WN:" + printPieces(b.WhitePieces.Knights) + "    " + "BN:" + printPieces(b.BlackPieces.Knights)
		case 2:
			result += "    WB:" + printPieces(b.WhitePieces.Bishops) + "    " + "BB:" + printPieces(b.BlackPieces.Bishops)
		case 3:
			result += "    WR:" + printPieces(b.WhitePieces.Rooks) + "    " + "BR:" + printPieces(b.BlackPieces.Rooks)
		case 4:
			result += "    WQ:" + printPieces(b.WhitePieces.Queens) + "    " + "BQ:" + printPieces(b.BlackPieces.Queens)
		case 5:
			result += "    WK:" + printPieces(b.WhitePieces.King) + "    " + "BK:" + printPieces(b.BlackPieces.King)

		}

		result += "\n"
	}
	return result
}

// GetLegalMoves generates all pseudo-legal moves for the current player. Does not check for uncovered king.
func (b *Board) GetLegalMoves() MoveList {
	var legalMoves MoveList

	for i := int8(0); i < 64; i++ {
		square := A1 << i
		row := i / 8
		col := i % 8
		pieceType, isWhite := b.PieceAtSquare(square)
		if pieceType == 0 || isWhite != b.WhiteToMove {
			continue
		}

		switch pieceType {
		case Pawn:
			if isWhite {
				// Generate pawn moves
				// Move forward
				if row < 7 {
					to := square << 8
					if !b.IsSquareOccupied(to) {
						move := b.NewMove(MoveNormal, square, to)
						move.EnPassant = 0
						legalMoves.Add(move)
					}
				}
				// Advance two squares from starting position
				if row == 1 {
					to := square << 16
					if !b.IsSquareOccupied(to) && !b.IsSquareOccupied(square<<8) {
						move := b.NewMove(MoveNormal, square, to)
						// En passant target = casilla intermedia
						move.EnPassant = uint8(bits.TrailingZeros64(square << 8))
						legalMoves.Add(move)
					}
				}
				// Capture diagonally to the right
				if row < 7 && col < 7 {
					to := square << 9
					if b.IsSquareOccupiedByBlack(to) {
						move := b.NewMove(MoveCapture, square, to)
						legalMoves.Add(move)
					} else if b.EnPassant != 0 && to == (uint64(1)<<b.EnPassant) { // en passant derecha
						move := b.NewMove(MoveCapture, square, to)
						legalMoves.Add(move)
					}
				}
				// Capture diagonally to the left
				if row < 7 && col > 0 {
					to := square << 7
					if b.IsSquareOccupiedByBlack(to) {
						move := b.NewMove(MoveCapture, square, to)
						legalMoves.Add(move)
					} else if b.EnPassant != 0 && to == (uint64(1)<<b.EnPassant) { // en passant izquierda
						move := b.NewMove(MoveCapture, square, to)
						legalMoves.Add(move)
					}
				}
			} else {
				// Generate pawn moves for black
				// Move forward
				if row > 0 {
					to := square >> 8
					if !b.IsSquareOccupied(to) {
						move := b.NewMove(MoveNormal, square, to)
						move.EnPassant = 0
						legalMoves.Add(move)
					}
				}
				// Advance two squares from starting position
				if row == 6 {
					to := square >> 16
					if !b.IsSquareOccupied(to) && !b.IsSquareOccupied(square>>8) {
						move := b.NewMove(MoveNormal, square, to)
						move.EnPassant = uint8(bits.TrailingZeros64(square >> 8))
						legalMoves.Add(move)
					}
				}
				// Capture diagonally to the right
				if row > 0 && col < 7 {
					to := square >> 7
					if b.IsSquareOccupiedByWhite(to) {
						move := b.NewMove(MoveCapture, square, to)
						legalMoves.Add(move)
					} else if b.EnPassant != 0 && to == (uint64(1)<<b.EnPassant) { // en passant derecha (desde negras)
						move := b.NewMove(MoveCapture, square, to)
						legalMoves.Add(move)
					}
				}
				// Capture diagonally to the left
				if row > 0 && col > 0 {
					to := square >> 9
					if b.IsSquareOccupiedByWhite(to) {
						move := b.NewMove(MoveCapture, square, to)
						legalMoves.Add(move)
					} else if b.EnPassant != 0 && to == (uint64(1)<<b.EnPassant) { // en passant izquierda (desde negras)
						move := b.NewMove(MoveCapture, square, to)
						legalMoves.Add(move)
					}
				}

			}
		case Knight:
			// Generate knight moves
			knightMoves := []int8{17, 15, 10, 6, -17, -15, -10, -6}
			for _, moveOffset := range knightMoves {
				toIndex := i + moveOffset
				toCol := toIndex % 8
				colDelta := math.Abs(float64(toCol - col)) // Ensure movement is not out of bounds
				if toIndex >= 0 && toIndex < 64 && (colDelta <= 2) {
					to := uint64(1) << toIndex
					occupiedByWhite := b.IsSquareOccupiedByWhite(to)
					occupiedByBlack := b.IsSquareOccupiedByBlack(to)
					// Knight cannot move to a square occupied by a piece of the same color
					moveForbidden := (isWhite && occupiedByWhite) || (!isWhite && occupiedByBlack)
					if !moveForbidden {
						// move := b.Clone()
						move := b.NewMove(MoveNormal, square, to)
						if (isWhite && occupiedByBlack) || (!isWhite && occupiedByWhite) {
							// move.CapturePiece(to, !isWhite)
							move.Type = MoveCapture
						}
						// move.MovePiece(square, to, isWhite, Knight)
						legalMoves.Add(move)
					}
				}
			}
		case Bishop:
			// Movimientos en las 4 diagonales: NE, NO, SE, SO
			for _, dir := range dirDiagonal {
				dirMoves := getLegalMovesInOneDirection(b, row, col, square, dir, isWhite)
				legalMoves.Append(dirMoves)
			}
		case Rook:
			rookMoves := []Move{}
			// Movimientos en las 4 direcciones: N, S, E, O
			for _, dir := range dirStraight {
				dirMoves := getLegalMovesInOneDirection(b, row, col, square, dir, isWhite)
				rookMoves = append(rookMoves, dirMoves...)
			}

			if isWhite {
				// Update white castling rights
				if square == A1 && (b.Castling&WhiteQueenSide != 0) {
					// Update castling rights for rook moves from original squares only if rights are set
					for i := range rookMoves {
						rookMoves[i].Castling &^= WhiteQueenSide
					}
				} else if square == H1 && (b.Castling&WhiteKingSide != 0) {
					for i := range rookMoves {
						rookMoves[i].Castling &^= WhiteKingSide
					}
				}
			} else {
				// Update black castling rights
				if square == A8 && (b.Castling&BlackQueenSide != 0) {
					// Update castling rights for rook moves from original squares only if rights are set
					for i := range rookMoves {
						rookMoves[i].Castling &^= BlackQueenSide
					}
				} else if square == H8 && (b.Castling&BlackKingSide != 0) {
					for i := range rookMoves {
						rookMoves[i].Castling &^= BlackKingSide
					}
				}
			}

			legalMoves = append(legalMoves, rookMoves...)
		case Queen:
			// Movimientos en las 8 direcciones: N, S, E, O, NE, NO, SE, SO
			for _, dir := range dirAll {
				dirMoves := getLegalMovesInOneDirection(b, row, col, square, dir, isWhite)
				legalMoves = append(legalMoves, dirMoves...)
			}
		case King:
			// Movimientos en las 8 direcciones pero solo una casilla
			for _, dir := range dirAll {
				r := row + dir.dr
				c := col + dir.dc
				if r >= 0 && r < 8 && c >= 0 && c < 8 {
					toIndex := r*8 + c
					to := uint64(1) << toIndex
					occupiedByWhite := b.IsSquareOccupiedByWhite(to)
					occupiedByBlack := b.IsSquareOccupiedByBlack(to)
					// King cannot move to a square occupied by a piece of the same color
					moveForbidden := ((isWhite && occupiedByWhite) || (!isWhite && occupiedByBlack)) ||
						b.SquareAttacked(r, c, isWhite)
					if !moveForbidden {
						// move := b.Clone()
						move := b.NewMove(MoveNormal, square, to)
						// Update castling rights
						if isWhite {
							move.Castling &^= WhiteKingSide | WhiteQueenSide
						} else {
							move.Castling &^= BlackKingSide | BlackQueenSide
						}
						if (isWhite && occupiedByBlack) || (!isWhite && occupiedByWhite) {
							// move.CapturePiece(to, !isWhite)
							move.Type = MoveCapture
						}
						// move.MovePiece(square, to, isWhite, King)
						legalMoves = append(legalMoves, move)
					}
				}
			}

			if square == E1 && isWhite {
				// Now check for castling rights
				if isWhite {
					// White king-side castling
					if b.Castling&WhiteKingSide != 0 &&
						!b.IsSquareOccupied(F1|G1) &&
						!b.SquareAttacked(0, 4, true) &&
						!b.SquareAttacked(0, 5, true) &&
						!b.SquareAttacked(0, 6, true) {
						move := b.NewMove(MoveKingCastle, square, G1)
						legalMoves = append(legalMoves, move)
					}
					// White queen-side castling
					if b.Castling&WhiteQueenSide != 0 &&
						!b.IsSquareOccupied(B1|C1|D1) &&
						!b.SquareAttacked(0, 4, true) &&
						!b.SquareAttacked(0, 3, true) &&
						!b.SquareAttacked(0, 2, true) {
						move := b.NewMove(MoveQueenCastle, square, C1)
						legalMoves = append(legalMoves, move)
					}
				}
			}
			if square == E8 && !isWhite {
				// Black king-side castling
				if b.Castling&BlackKingSide != 0 &&
					!b.IsSquareOccupied(F8|G8) &&
					!b.SquareAttacked(7, 4, false) &&
					!b.SquareAttacked(7, 5, false) &&
					!b.SquareAttacked(7, 6, false) {
					move := b.NewMove(MoveKingCastle, square, G8)
					legalMoves = append(legalMoves, move)
				}
				// Black queen-side castling
				if b.Castling&BlackQueenSide != 0 &&
					!b.IsSquareOccupied(B8|C8|D8) &&
					!b.SquareAttacked(7, 4, false) &&
					!b.SquareAttacked(7, 3, false) &&
					!b.SquareAttacked(7, 2, false) {
					move := b.NewMove(MoveQueenCastle, square, C8)
					legalMoves = append(legalMoves, move)
				}

			}
		default:
			// For simplicity, other pieces are not implemented in this example
		}
	}
	return legalMoves
}

type Direction struct{ dr, dc int8 }

var dirStraight = []Direction{
	{1, 0},  // N
	{-1, 0}, // S
	{0, 1},  // E
	{0, -1}, // O
}
var dirDiagonal = []Direction{
	{1, 1},   // NE
	{1, -1},  // NO
	{-1, 1},  // SE
	{-1, -1}, // SO
}
var dirAll = []Direction{
	{1, 0},   // N
	{-1, 0},  // S
	{0, 1},   // E
	{0, -1},  // O
	{1, 1},   // NE
	{1, -1},  // NO
	{-1, 1},  // SE
	{-1, -1}, // SO
}

// getLegalMovesInOneDirection generates legal moves for sliding pieces (Bishop, Rook, Queen) in a given direction.
func getLegalMovesInOneDirection(b *Board, r int8, c int8, square uint64, dir Direction, isWhite bool) []Move {
	var legalMoves []Move
	// from := r*8 + c
	for {
		r += dir.dr
		c += dir.dc
		if r < 0 || r > 7 || c < 0 || c > 7 {
			break
		}
		toIndex := r*8 + c
		to := uint64(1) << toIndex
		if b.IsSquareOccupied(to) {
			blackCapture := isWhite && b.IsSquareOccupiedByBlack(to)
			whiteCapture := !isWhite && b.IsSquareOccupiedByWhite(to)
			// Si es pieza enemiga, permite captura
			if whiteCapture || blackCapture {
				// move := b.Clone()
				// move.CapturePiece(to, whiteCapture)
				// move.MovePiece(square, to, isWhite, piece)
				move := b.NewMove(MoveCapture, square, to)
				legalMoves = append(legalMoves, move)
			}
			break // No puede saltar piezas
		} else {
			// move := b.Clone()
			// move.MovePiece(square, to, isWhite, piece)
			move := b.NewMove(MoveNormal, square, to)
			legalMoves = append(legalMoves, move)
		}
	}
	return legalMoves
}

// SquareAttacked checks if the given square is attacked by any piece of the opponent.
// isWhite is the color of the player whose king is being checked (i.e., isWhite=true means check if black attacks).
// row and col are 0-indexed (0-7).
func (b *Board) SquareAttacked(row, col int8, isWhite bool) bool {
	// fmt.Println("Checking attack on square:", row, col, "isWhite:", isWhite)
	// if row == 0 && col == 4 && !isWhite {
	// 	fmt.Println("Check bug")
	// 	fmt.Println(b.ToString())
	// }
	var opponentPieces *Pieces
	if isWhite {
		opponentPieces = &b.BlackPieces
	} else {
		opponentPieces = &b.WhitePieces
	}

	getBitSquare := func(r, c int8) uint64 {
		return uint64(1) << (r*8 + c)
	}
	pawns := opponentPieces.Pawns
	// Pawns that can attack the square
	if isWhite && row < 7 {
		pawns = getBitSquare(row+1, col-1) | getBitSquare(row+1, col+1)
		if pawns&opponentPieces.Pawns != 0 {
			return true
		}
	} else if !isWhite && row > 0 {
		pawns = getBitSquare(row-1, col-1) | getBitSquare(row-1, col+1)
		if pawns&opponentPieces.Pawns != 0 {
			return true
		}
	}

	// Knights that can attack the square
	knightMoves := []struct{ dr, dc int8 }{
		{2, 1}, {2, -1}, {1, 2}, {1, -2},
		{-2, 1}, {-2, -1}, {-1, 2}, {-1, -2},
	}
	for _, move := range knightMoves {
		r := row + move.dr
		c := col + move.dc
		if r >= 0 && r < 8 && c >= 0 && c < 8 {
			if opponentPieces.Knights&getBitSquare(r, c) != 0 {
				return true
			}
		}
	}

	// Bishops/Queens that can attack the square diagonally
	for _, dir := range dirDiagonal {
		rTemp, cTemp := row, col
		for {
			rTemp += dir.dr
			cTemp += dir.dc
			if rTemp < 0 || rTemp > 7 || cTemp < 0 || cTemp > 7 {
				break
			}
			to := getBitSquare(rTemp, cTemp)
			if b.IsSquareOccupied(to) {
				if opponentPieces.Bishops&to != 0 || opponentPieces.Queens&to != 0 {
					return true
				}
				break // No puede saltar piezas
			}
		}
	}

	// Rooks/Queens that can attack the square orthogonally
	for _, dir := range dirStraight {
		rTemp, cTemp := row, col
		for {
			rTemp += dir.dr
			cTemp += dir.dc
			if rTemp < 0 || rTemp > 7 || cTemp < 0 || cTemp > 7 {
				break
			}
			to := getBitSquare(rTemp, cTemp)
			if b.IsSquareOccupied(to) {
				if opponentPieces.Rooks&to != 0 || opponentPieces.Queens&to != 0 {
					return true
				}
				break // No puede saltar piezas
			}
		}
	}

	// Opponent King that can attack the square
	for _, dir := range dirAll {
		rTemp := row + dir.dr
		cTemp := col + dir.dc
		if rTemp >= 0 && rTemp < 8 && cTemp >= 0 && cTemp < 8 {
			to := getBitSquare(rTemp, cTemp)
			if opponentPieces.King&to != 0 {
				return true
			}
		}
	}

	return false
}
