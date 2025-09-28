package main

import (
	"fmt"
	m "zentense/muabdib"
)

func main() {

	board := m.NewBoard()
	board.SetFen("r3k1nr/ppp2ppp/8/4p1B1/3Nq3/6P1/PP1QBP1P/R3K2R b KQkq - 0 1")
	legalMoves := board.GetLegalMoves()
	fmt.Println(legalMoves.ToString(true))
	// m, bestScore := board.GetBestMove(2, true)
	// fmt.Println(m.ToString(), bestScore)
	mseq, bestScore := board.GetBestLine(5, false)
	fmt.Println(mseq.ToString(false), bestScore)
	return

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

}
