package melange

import (
	"math/bits"
	"math/rand"
)

// Zobrist holds the random numbers for Zobrist hashing.
var Zobrist struct {
	// pieceKeys[color][pieceType][square]
	// color: 0 for White, 1 for Black
	// pieceType: 0 for Pawn, 1 for Knight, ..., 5 for King (Piece - 1)
	pieceKeys [2][6][64]uint64
	// castlingKeys[castlingRights]
	castlingKeys [16]uint64
	// enPassantKeys[file] - only 8 files are possible for en passant
	enPassantKeys [8]uint64
	// sideToMoveKey is XORed if it's black's turn.
	sideToMoveKey uint64
}

// init initializes the Zobrist keys with random values.
func init() {
	// Seed the random number generator.
	// Using a fixed seed for deterministic keys.
	random := rand.New(rand.NewSource(3141592653589793238))

	for i := 0; i < 2; i++ {
		for j := 0; j < 6; j++ {
			for k := 0; k < 64; k++ {
				Zobrist.pieceKeys[i][j][k] = random.Uint64()
			}
		}
	}

	for i := 0; i < 16; i++ {
		Zobrist.castlingKeys[i] = random.Uint64()
	}

	for i := 0; i < 8; i++ {
		Zobrist.enPassantKeys[i] = random.Uint64()
	}

	Zobrist.sideToMoveKey = random.Uint64()
}

func (b *Board) CalculateHash() uint64 {
	var hash uint64

	// Helper to hash a bitboard
	hashBitboard := func(bb uint64, colorIndex int, pieceIdx int) {
		for bb != 0 {
			sq := bits.TrailingZeros64(bb)
			hash ^= Zobrist.pieceKeys[colorIndex][pieceIdx][sq]
			bb &= bb - 1
		}
	}

	// White pieces (colorIndex 0)
	hashBitboard(b.WhitePieces.Pawns, 0, int(Pawn)-1)
	hashBitboard(b.WhitePieces.Knights, 0, int(Knight)-1)
	hashBitboard(b.WhitePieces.Bishops, 0, int(Bishop)-1)
	hashBitboard(b.WhitePieces.Rooks, 0, int(Rook)-1)
	hashBitboard(b.WhitePieces.Queens, 0, int(Queen)-1)
	hashBitboard(b.WhitePieces.King, 0, int(King)-1)

	// Black pieces (colorIndex 1)
	hashBitboard(b.BlackPieces.Pawns, 1, int(Pawn)-1)
	hashBitboard(b.BlackPieces.Knights, 1, int(Knight)-1)
	hashBitboard(b.BlackPieces.Bishops, 1, int(Bishop)-1)
	hashBitboard(b.BlackPieces.Rooks, 1, int(Rook)-1)
	hashBitboard(b.BlackPieces.Queens, 1, int(Queen)-1)
	hashBitboard(b.BlackPieces.King, 1, int(King)-1)

	// Castling rights
	hash ^= Zobrist.castlingKeys[b.Castling]

	// En passant: only if there is a pawn that can capture hash must be different
	if b.EnPassant != 0 {
		hasAttacker := false
		epBit := uint64(1) << b.EnPassant
		col := b.EnPassant % 8

		if b.WhiteToMove {
			// Check for white pawns that can capture the EP square
			possibleAttackers := uint64(0)
			if col > 0 {
				possibleAttackers |= (epBit >> 9) // Attack from left (file-1)
			}
			if col < 7 {
				possibleAttackers |= (epBit >> 7) // Attack from right (file+1)
			}
			if (possibleAttackers & b.WhitePieces.Pawns) != 0 {
				hasAttacker = true
			}
		} else {
			// Check for black pawns that can capture the EP square
			possibleAttackers := uint64(0)
			if col > 0 {
				possibleAttackers |= (epBit << 7) // Attack from left (file-1)
			}
			if col < 7 {
				possibleAttackers |= (epBit << 9) // Attack from right (file+1)
			}
			if (possibleAttackers & b.BlackPieces.Pawns) != 0 {
				hasAttacker = true
			}
		}

		if hasAttacker {
			file := (b.EnPassant) % 8
			hash ^= Zobrist.enPassantKeys[file]
		}
	}

	// Side to move
	if !b.WhiteToMove {
		hash ^= Zobrist.sideToMoveKey
	}

	return hash
}
