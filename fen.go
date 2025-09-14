package muabdib

import "fmt"

// SetFEN inicializa el tablero a partir de un string FEN
func (b *Board) SetFen(fen string) error {
	// Reset all pieces
	b.WhitePieces = Pieces{}
	b.BlackPieces = Pieces{}

	parts := make([]string, 6)
	n := copy(parts, splitFEN(fen))
	if n < 4 {
		return fmt.Errorf("FEN inválido: %s", fen)
	}

	// Parse piezas
	rank := 7
	file := 0
	for _, ch := range parts[0] {
		if ch == '/' {
			rank--
			file = 0
			continue
		}
		if ch >= '1' && ch <= '8' {
			file += int(ch - '0')
			continue
		}
		sq := uint64(1) << (rank*8 + file)
		switch ch {
		case 'P':
			b.WhitePieces.Pawns |= sq
		case 'N':
			b.WhitePieces.Knights |= sq
		case 'B':
			b.WhitePieces.Bishops |= sq
		case 'R':
			b.WhitePieces.Rooks |= sq
		case 'Q':
			b.WhitePieces.Queens |= sq
		case 'K':
			b.WhitePieces.King |= sq
		case 'p':
			b.BlackPieces.Pawns |= sq
		case 'n':
			b.BlackPieces.Knights |= sq
		case 'b':
			b.BlackPieces.Bishops |= sq
		case 'r':
			b.BlackPieces.Rooks |= sq
		case 'q':
			b.BlackPieces.Queens |= sq
		case 'k':
			b.BlackPieces.King |= sq
		default:
			return fmt.Errorf("FEN inválido: pieza desconocida '%c'", ch)
		}
		file++
	}

	// Parse turno
	b.WhiteToMove = parts[1] == "w"

	// Parse derechos de enroque
	b.Castling = 0
	for _, ch := range parts[2] {
		switch ch {
		case 'K':
			b.Castling |= WhiteKingSide
		case 'Q':
			b.Castling |= WhiteQueenSide
		case 'k':
			b.Castling |= BlackKingSide
		case 'q':
			b.Castling |= BlackQueenSide
		case '-':
		default:
			return fmt.Errorf("FEN inválido: derecho de enroque desconocido '%c'", ch)
		}
	}

	// Parse en passant
	b.EnPassant = 0
	if parts[3] != "-" {
		epSq, err := parseFENSquare(parts[3])
		if err != nil {
			return err
		}
		b.EnPassant = uint8(epSq)
	}

	// Parse halfmove y fullmove si existen
	if n > 4 {
		fmt.Sscanf(parts[4], "%d", &b.HalfMove)
	} else {
		b.HalfMove = 0
	}
	if n > 5 {
		fmt.Sscanf(parts[5], "%d", &b.FullMove)
	} else {
		b.FullMove = 1
	}
	return nil
}

// splitFEN separa el string FEN en sus partes
func splitFEN(fen string) []string {
	out := make([]string, 0, 6)
	curr := ""
	for _, ch := range fen {
		if ch == ' ' {
			out = append(out, curr)
			curr = ""
		} else {
			curr += string(ch)
		}
	}
	out = append(out, curr)
	return out
}

// parseFENSquare convierte una notación algebraica (ej. "e3") a índice 0-63
func parseFENSquare(sq string) (int, error) {
	if len(sq) != 2 {
		return 0, fmt.Errorf("FEN inválido: en passant '%s'", sq)
	}
	file := int(sq[0] - 'a')
	rank := int(sq[1] - '1')
	if file < 0 || file > 7 || rank < 0 || rank > 7 {
		return 0, fmt.Errorf("FEN inválido: en passant '%s'", sq)
	}
	return rank*8 + file, nil
}
