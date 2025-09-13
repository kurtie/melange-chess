package muabdib

import (
	"fmt"
	"sort"
	"strings"
)

const (
	A1 uint64 = 1 << 0
	B1 uint64 = 1 << 1
	C1 uint64 = 1 << 2
	D1 uint64 = 1 << 3
	E1 uint64 = 1 << 4
	F1 uint64 = 1 << 5
	G1 uint64 = 1 << 6
	H1 uint64 = 1 << 7
	A2 uint64 = 1 << 8
	B2 uint64 = 1 << 9
	C2 uint64 = 1 << 10
	D2 uint64 = 1 << 11
	E2 uint64 = 1 << 12
	F2 uint64 = 1 << 13
	G2 uint64 = 1 << 14
	H2 uint64 = 1 << 15
	A3 uint64 = 1 << 16
	B3 uint64 = 1 << 17
	C3 uint64 = 1 << 18
	D3 uint64 = 1 << 19
	E3 uint64 = 1 << 20
	F3 uint64 = 1 << 21
	G3 uint64 = 1 << 22
	H3 uint64 = 1 << 23
	A4 uint64 = 1 << 24
	B4 uint64 = 1 << 25
	C4 uint64 = 1 << 26
	D4 uint64 = 1 << 27
	E4 uint64 = 1 << 28
	F4 uint64 = 1 << 29
	G4 uint64 = 1 << 30
	H4 uint64 = 1 << 31
	A5 uint64 = 1 << 32
	B5 uint64 = 1 << 33
	C5 uint64 = 1 << 34
	D5 uint64 = 1 << 35
	E5 uint64 = 1 << 36
	F5 uint64 = 1 << 37
	G5 uint64 = 1 << 38
	H5 uint64 = 1 << 39
	A6 uint64 = 1 << 40
	B6 uint64 = 1 << 41
	C6 uint64 = 1 << 42
	D6 uint64 = 1 << 43
	E6 uint64 = 1 << 44
	F6 uint64 = 1 << 45
	G6 uint64 = 1 << 46
	H6 uint64 = 1 << 47
	A7 uint64 = 1 << 48
	B7 uint64 = 1 << 49
	C7 uint64 = 1 << 50
	D7 uint64 = 1 << 51
	E7 uint64 = 1 << 52
	F7 uint64 = 1 << 53
	G7 uint64 = 1 << 54
	H7 uint64 = 1 << 55
	A8 uint64 = 1 << 56
	B8 uint64 = 1 << 57
	C8 uint64 = 1 << 58
	D8 uint64 = 1 << 59
	E8 uint64 = 1 << 60
	F8 uint64 = 1 << 61
	G8 uint64 = 1 << 62
	H8 uint64 = 1 << 63
)

type Piece int

const (
	Pawn   Piece = 1
	Knight Piece = 2
	Bishop Piece = 3
	Rook   Piece = 4
	Queen  Piece = 5
	King   Piece = 6
)

type Pieces struct {
	Pawns   uint64
	Knights uint64
	Bishops uint64
	Rooks   uint64
	Queens  uint64
	King    uint64
}

func (p *Pieces) Get(piece Piece) uint64 {
	switch piece {
	case Pawn:
		return p.Pawns
	case Knight:
		return p.Knights
	case Bishop:
		return p.Bishops
	case Rook:
		return p.Rooks
	case Queen:
		return p.Queens
	case King:
		return p.King
	default:
		return 0
	}
}

type CastleRights uint8

const (
	WhiteKingSide  CastleRights = 1 << 0
	WhiteQueenSide CastleRights = 1 << 1
	BlackKingSide  CastleRights = 1 << 2
	BlackQueenSide CastleRights = 1 << 3
)

type MoveType uint8

const (
	MoveNormal             MoveType = iota
	MoveCapture            MoveType = 1
	MoveKingCastle         MoveType = 2
	MoveQueenCastle        MoveType = 4
	MoveDoublePawn         MoveType = 6
	MovePromotion          MoveType = 8
	MoveKnightPromo        MoveType = MovePromotion | 16
	MoveBishopPromo        MoveType = MovePromotion | 32
	MoveRookPromo          MoveType = MovePromotion | 64
	MoveQueenPromo         MoveType = MovePromotion | 128
	MoveKnightPromoCapture MoveType = MoveKnightPromo | MoveCapture
	MoveBishopPromoCapture MoveType = MoveBishopPromo | MoveCapture
	MoveRookPromoCapture   MoveType = MoveRookPromo | MoveCapture
	QueenPromoCapture      MoveType = MoveQueenPromo | MoveCapture
)

type Move struct {
	From      int8
	To        int8
	Type      MoveType
	Castling  CastleRights
	EnPassant uint8
}

func (m *Move) ToString() string {
	from := squareToString(m.From)
	to := squareToString(m.To)
	mt := m.Type
	// Enroques
	if mt == MoveKingCastle {
		return "O-O"
	}
	if mt == MoveQueenCastle {
		return "O-O-O"
	}
	// Promociones con captura
	if (mt&MoveCapture) != 0 && (mt&MovePromotion) != 0 {
		promo := ""
		if (mt & 16) != 0 {
			promo = "=N"
		} else if (mt & 32) != 0 {
			promo = "=B"
		} else if (mt & 64) != 0 {
			promo = "=R"
		} else if (mt & 128) != 0 {
			promo = "=Q"
		}
		return fmt.Sprintf("%sx%s%s%s", from, to, promo, "")
	}
	// Promociones sin captura
	if (mt & MovePromotion) != 0 {
		promo := ""
		if (mt & 16) != 0 {
			promo = "=N"
		} else if (mt & 32) != 0 {
			promo = "=B"
		} else if (mt & 64) != 0 {
			promo = "=R"
		} else if (mt & 128) != 0 {
			promo = "=Q"
		}
		return fmt.Sprintf("%s%s%s", from, to, promo)
	}
	// Captura normal
	if (mt & MoveCapture) != 0 {
		return fmt.Sprintf("%sx%s", from, to)
	}
	// Movimiento normal
	return fmt.Sprintf("%s%s", from, to)
}

// squareToString convierte un índice de 0-63 a notación tipo A1, C3, etc.
func squareToString(idx int8) string {
	file := idx % 8
	rank := (idx / 8)
	return fmt.Sprintf("%c%d", 'a'+file, rank+1)
}

type MoveList []Move

func (ml *MoveList) Add(move Move) {
	*ml = append(*ml, move)
}

// Append añade todos los movimientos de otro MoveList a este MoveList
func (ml *MoveList) Append(other MoveList) {
	*ml = append(*ml, other...)
}

func (ml *MoveList) ToStringArray() []string {
	arr := make([]string, len(*ml))
	for i, m := range *ml {
		arr[i] = m.ToString()
	}
	return arr
}

func (ml *MoveList) ToStringArraySorted() []string {
	arr := ml.ToStringArray()
	sort.Strings(arr)
	return arr
}

func (ml *MoveList) ToString(sorted bool) string {
	if sorted {
		return strings.Join(ml.ToStringArraySorted(), ", ")
	}
	return strings.Join(ml.ToStringArray(), ", ")
}
