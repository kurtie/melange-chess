package main

import (
	"fmt"
	"zentense/muabdib"
)

func main() {

	board := muabdib.NewBoard()
	// board.WhiteToMove = false
	fmt.Println(board.ToString())
	return
	// legalMoves := board.GetLegalMoves()
	// for _, move := range legalMoves {
	// 	fmt.Println(move.ToString())
	// }
	// fmt.Println("Total legal moves:", len(legalMoves))

	board = &muabdib.Board{}
	board.WhiteToMove = true
	board.WhitePieces.King = (uint64(1) << 36)  // King on g5
	board.BlackPieces.Rooks = (uint64(1) << 37) // Bishop on g6
	board.WhitePieces.Pawns = (uint64(1) << 35) // Pawn on g4
	fmt.Println(board.ToString())
	legalMoves := board.GetLegalMoves()
	for _, move := range legalMoves {
		fmt.Println(move.ToString())
	}
	fmt.Println("Total legal moves:", len(legalMoves))

	// board = &muabdib.Board{
	// 	WhitePieces: muabdib.Pieces{
	// 		Pawns: uint64(0x0000000800000000),
	// 		King:  uint64(0x0000001000000000),
	// 	},
	// 	BlackPieces: muabdib.Pieces{
	// 		Pawns: uint64(0x0000002000000000),
	// 	},
	// 	WhiteToMove: true,
	// }
	// println(board.ToString())
	// println("Attacked E4:", board.SquareAttacked(3, 4, false))
	// println("Attacked E4:", board.SquareAttacked(3, 5, false))
	// println("Attacked E4:", board.SquareAttacked(3, 6, false))
}
