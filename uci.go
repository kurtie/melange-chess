package muabdib

import (
	"fmt"
)

// currentBoard holds the persistent board state across UCI commands
var currentBoard *Board

// GetCurrentBoard returns the board managed by the UCI interface (for tests/inspection)
func GetCurrentBoard() *Board {
	if currentBoard == nil {
		currentBoard = NewBoard()
	}
	return currentBoard
}

func ProcessUciCommand(command string) bool {
	fmt.Println("debug: received command:", command)
	// Tokenize command string, allowing arbitrary whitespace between tokens
	// Remove trailing newline
	trimmed := command
	if len(trimmed) > 0 && trimmed[len(trimmed)-1] == '\n' {
		trimmed = trimmed[:len(trimmed)-1]
	}
	tokens := tokenize(trimmed)
	if len(tokens) > 0 {
		switch tokens[0] {
		case "go":
			handleGo(tokens)
		case "isready":
			// Initializations done here
			fmt.Println("readyok")
		case "position":
			handlePosition(tokens)
		case "quit":
			fmt.Println("Exiting...")
			return true
		case "uci":
			// Use fmt.Println for proper newline handling
			fmt.Println("id name Muabdib v0.1")
			fmt.Println("id author Jose R. Cabanes")
			fmt.Println("uciok")
		case "ucinewgame":
			// Reset engine state for a new game
			currentBoard = NewBoard()
		default:
			fmt.Println("Unknown command:", command)
		}
	}
	return false
}

// tokenize splits a string into tokens separated by arbitrary whitespace
func tokenize(s string) []string {
	tokens := []string{}
	token := ""
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == ' ' || c == '\t' {
			if token != "" {
				tokens = append(tokens, token)
				token = ""
			}
			continue
		}
		token += string(c)
	}
	if token != "" {
		tokens = append(tokens, token)
	}
	return tokens
}

// handlePosition parses and applies the UCI 'position' command
// Syntax: position [fen <fenstring> | startpos ] [moves <move1> .... <movei>]
func handlePosition(tokens []string) {
	if len(tokens) < 2 {
		fmt.Println("info string invalid position command", tokens)
		return
	}

	// Ensure board exists
	if currentBoard == nil {
		currentBoard = NewBoard()
	}

	idx := 1
	// Setup position base
	if tokens[idx] == "startpos" {
		currentBoard = NewBoard()
		idx++
	} else if tokens[idx] == "fen" {
		idx++
		// Collect fen parts until 'moves' or end
		fenStart := idx
		for idx < len(tokens) && tokens[idx] != "moves" {
			idx++
		}
		fen := joinWithSpaces(tokens[fenStart:idx])
		b := &Board{}
		if err := b.SetFen(fen); err != nil {
			// On invalid FEN, keep previous board but report
			fmt.Println("info string invalid FEN:", err)
			return
		}
		currentBoard = b
	} else {
		// Unknown base, ignore
		return
	}

	// Apply moves if present
	if idx < len(tokens) && tokens[idx] == "moves" {
		idx++
		for idx < len(tokens) {
			mvStr := tokens[idx]
			mv, ok := parseUCIMove(currentBoard, mvStr)
			if !ok {
				// If a move cannot be parsed/applied, report and stop applying further
				fmt.Println("info string invalid move:", mvStr)
				return
			}
			// Apply move
			fromBB := mv.GetFrom64()
			_, isWhite := currentBoard.PieceAtSquare(fromBB)
			currentBoard.MovePiece(mv, isWhite)
			// Update clocks (best-effort): halfmove resets on pawn move or capture; fullmove after Black's move
			if mv.Piece == Pawn || mv.IsCapture() {
				currentBoard.HalfMove = 0
			} else {
				currentBoard.HalfMove++
			}
			if currentBoard.WhiteToMove { // after toggled in MovePiece; if now White to move, Black just moved
				currentBoard.FullMove++
			}
			idx++
		}
	}
}

func handleGo(tokens []string) {
	// Handle the 'go' command
	if len(tokens) < 2 {
		fmt.Println("info string invalid go command", tokens)
		return
	}
	if currentBoard == nil {
		currentBoard = NewBoard()
	}
	// Start the search for the best move
	root := currentBoard.GenerateSearchTree(5)
	move := root.GetBestMove()
	// fmt.Println(currentBoard.ToString())
	fmt.Println("bestmove", move.ToSimpleString())
}

func joinWithSpaces(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	// manual join to avoid importing strings just for this
	out := parts[0]
	for i := 1; i < len(parts); i++ {
		out += " " + parts[i]
	}
	return out
}

// parseUCIMove finds and returns the legal move matching the UCI long algebraic string (e2e4[,qrbn])
func parseUCIMove(b *Board, uci string) (Move, bool) {
	// Expect 4 or 5 chars
	if len(uci) < 4 {
		return Move{}, false
	}
	fromStr := uci[0:2]
	toStr := uci[2:4]
	promo := byte(0)
	if len(uci) >= 5 {
		promo = uci[4]
	}
	fromIdx, ok := parseSquareToIndex(fromStr)
	if !ok {
		return Move{}, false
	}
	toIdx, ok := parseSquareToIndex(toStr)
	if !ok {
		return Move{}, false
	}

	moves := b.GetLegalMoves()
	for _, m := range moves {
		// Filter moves that are illegal (leave king in check)
		if !b.isMoveLegal(m) {
			continue
		}
		if int(m.From) != fromIdx || int(m.To) != toIdx {
			continue
		}
		// If promotion present, ensure types match; if not present, skip promotion moves
		isPromo := (m.Type & MovePromotion) != 0
		if promo == 0 && isPromo {
			continue
		}
		if promo != 0 {
			switch promo {
			case 'q':
				if (m.Type & 128) == 0 { // Queen flag in our encoding
					continue
				}
			case 'r':
				if (m.Type & 64) == 0 {
					continue
				}
			case 'b':
				if (m.Type & 32) == 0 {
					continue
				}
			case 'n':
				if (m.Type & 16) == 0 {
					continue
				}
			default:
				continue
			}
		}
		return m, true
	}
	return Move{}, false
}

// parseSquareToIndex converts algebraic square like "e2" to index 0..63
func parseSquareToIndex(s string) (int, bool) {
	if len(s) != 2 {
		return 0, false
	}
	f := int(s[0] - 'a')
	r := int(s[1] - '1')
	if f < 0 || f > 7 || r < 0 || r > 7 {
		return 0, false
	}
	return r*8 + f, true
}
