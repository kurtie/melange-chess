package main

import (
	"fmt"
	m "zentense/muabdib"
)

func main() {

	board := m.NewBoard()
	// board.WhiteToMove = true
	// board.MovePiece(board.NewMove(m.MoveNormal, m.E2, m.E4), true, m.Pawn)
	// board.MovePiece(board.NewMove(m.MoveNormal, m.D7, m.D5), false, m.Pawn)
	// board.MovePiece(board.NewMove(m.MoveNormal, m.F1, m.B5), true, m.Bishop)
	board.SetFen("rnbqkbnr/ppp1pppp/8/1B1p4/4P3/8/PPPP1PPP/RNBQK1NR w KQkq - 0 1")
	fmt.Println(board.ToString())
	fmt.Println(board.SquareAttacked(7, 4, true)) // e1

	board.SetFen("8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1")
	fmt.Println(board.ToString())
	return

	res := m.Perft(board, 2)
	fmt.Println("Perft results at depth 2:")
	fmt.Printf("Nodes: %d\n", res.Nodes)
	fmt.Printf("Captures: %d\n", res.Captures)
	fmt.Printf("En Passants: %d\n", res.EnPassants)
	fmt.Printf("Checks: %d\n", res.Checks)
}
