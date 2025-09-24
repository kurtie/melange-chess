package muabdib

import "math/bits"

type Centipawn int

const (
	CpPawn         int = 100
	CpKnight       int = 320
	CpBishop       int = 330
	CpRook         int = 500
	CpQueen        int = 900
	CpKing         int = 10000
	CpDoubledPawn  int = 50
	CpStuckedPawn  int = 30
	CpIsolatedPawn int = 20
	CpBishopPair   int = 30
)

// Posicionamiento de piezas (piece-square tables)
// Valores en centipawns
// Las tablas están definidas desde la perspectiva de las blancas visualmente.
// Para las blancas, se debe invertir el índice (63 - sq).
var PosPawn = []int{
	0, 0, 0, 0, 0, 0, 0, 0,
	50, 50, 50, 50, 50, 50, 50, 50,
	10, 10, 20, 30, 30, 20, 10, 10,
	5, 5, 10, 25, 25, 10, 5, 5,
	0, 0, 0, 20, 20, 0, 0, 0,
	5, -5, -10, 0, 0, -10, -5, 5,
	5, 10, 10, -20, -20, 10, 10, 5,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var PosKnight = []int{
	-50, -40, -30, -30, -30, -30, -40, -50,
	-40, -20, 0, 0, 0, 0, -20, -40,
	-30, 0, 10, 15, 15, 10, 0, -30,
	-30, 5, 15, 20, 20, 15, 5, -30,
	-30, 0, 15, 20, 20, 15, 0, -30,
	-30, 5, 10, 15, 15, 10, 5, -30,
	-40, -20, 0, 5, 5, 0, -20, -40,
	-50, -40, -30, -30, -30, -30, -40, -50,
}

var PosBishop = []int{
	-20, -10, -10, -10, -10, -10, -10, -20,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 5, 10, 10, 5, 0, -10,
	-10, 5, 5, 10, 10, 5, 5, -10,
	-10, 0, 10, 10, 10, 10, 0, -10,
	-10, 10, 10, 10, 10, 10, 10, -10,
	-10, 5, 0, 0, 0, 0, 5, -10,
	-20, -10, -10, -10, -10, -10, -10, -20,
}

var PosRook = []int{
	0, 0, 0, 0, 0, 0, 0, 0,
	5, 10, 10, 10, 10, 10, 10, 5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	0, 0, 0, 5, 5, 0, 0, 0,
}

var PosQueen = []int{
	-20, -10, -10, -5, -5, -10, -10, -20,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 5, 5, 5, 5, 0, -10,
	-5, 0, 5, 5, 5, 5, 0, -5,
	0, 0, 5, 5, 5, 5, 0, -5,
	-10, 5, 5, 5, 5, 5, 0, -10,
	-10, 0, 5, 0, 0, 0, 0, -10,
	-20, -10, -10, -5, -5, -10, -10, -20,
}

var PosKingMiddle = []int{
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-20, -30, -30, -40, -40, -30, -30, -20,
	-10, -20, -20, -20, -20, -20, -20, -10,
	20, 20, 0, 0, 0, 0, 20, 20,
	20, 30, 10, 0, 0, 10, 30, 20,
}

var PosKingEnd = []int{
	-50, -40, -30, -20, -20, -30, -40, -50,
	-30, -20, -10, 0, 0, -10, -20, -30,
	-30, -10, 20, 30, 30, 20, -10, -30,
	-30, -10, 30, 40, 40, 30, -10, -30,
	-30, -10, 30, 40, 40, 30, -10, -30,
	-30, -10, 20, 30, 30, 20, -10, -30,
	-30, -30, 0, 0, 0, 0, -30, -30,
	-50, -30, -30, -30, -30, -30, -30, -50,
}

func (b *Board) Evaluate() int {
	// Suma simple de material: (blancas - negras)
	// Usa popcount para cada bitboard y multiplica por el valor de la pieza.
	scoreWhite := b.WhitePieces.evalMaterial() + b.WhitePieces.evalPositions(true)
	scoreBlack := b.BlackPieces.evalMaterial() + b.BlackPieces.evalPositions(false)

	return scoreWhite - scoreBlack
}

func (p *Pieces) evalMaterial() int {
	material := 0

	material += bitsOnesCount(p.Pawns) * CpPawn
	material += bitsOnesCount(p.Knights) * CpKnight
	material += bitsOnesCount(p.Bishops) * CpBishop
	material += bitsOnesCount(p.Rooks) * CpRook
	material += bitsOnesCount(p.Queens) * CpQueen
	material += bitsOnesCount(p.King) * CpKing

	return material
}

// bitsOnesCount es una pequeña envoltura; separado para facilitar tests si se desea.
func bitsOnesCount(bb uint64) int {
	return bits.OnesCount64(bb)
}

func (p *Pieces) evalPositions(isWhite bool) int {
	score := 0

	// Helper to accumulate score from a bitboard and its piece-square table.
	accumulate := func(bb uint64, table []int) {
		for bb != 0 {
			sq := bits.TrailingZeros64(bb) // 0..63 (A1 = 0)
			bb &= bb - 1                   // clear lowest bit
			idx := sq
			if isWhite { // Las tablas están definidas desde la perspectiva visual de las blancas (A8 primero)
				// Para piezas blancas invertir el índice (63 - sq) según comentario en archivo.
				idx = 63 - sq
			}
			score += table[idx]
		}
	}

	accumulate(p.Pawns, PosPawn)
	accumulate(p.Knights, PosKnight)
	accumulate(p.Bishops, PosBishop)
	accumulate(p.Rooks, PosRook)
	accumulate(p.Queens, PosQueen)

	// Para el rey usamos siempre la tabla de medio juego por ahora.
	// (Una heurística de final podría añadirse más adelante detectando material reducido.)
	accumulate(p.King, PosKingMiddle)

	return score
}
