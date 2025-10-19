package melange

import (
	"fmt"
	"slices"
	"sort"
)

type SearchNode struct {
	Parent   *SearchNode
	Children []*SearchNode
	Move     Move
	Score    int
}

var totalNodes = 0

// internal slice type to avoid reflect overhead of sort.Slice
type searchNodeSlice []*SearchNode

func (s searchNodeSlice) Len() int      { return len(s) }
func (s searchNodeSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Descending order by Score (may place best moves first for alpha-beta move ordering)
func (s searchNodeSlice) Less(i, j int) bool { return s[i].Score > s[j].Score }

// sortChildren ordena los hijos por Score descendente usando sort.Sort sobre
// un tipo especializado para evitar el coste de reflexión de sort.Slice.
// Para tamaños pequeños (típico branching factor del ajedrez < 40) el algoritmo
// interno (quicksort con inserción) es muy eficiente.
func (n *SearchNode) sortChildren(whiteToMove bool) {
	if len(n.Children) < 2 {
		return
	}
	if whiteToMove {
		slices.SortFunc(n.Children, func(a, b *SearchNode) int {
			return b.Score - a.Score
		})
	} else {
		slices.SortFunc(n.Children, func(a, b *SearchNode) int {
			return a.Score - b.Score
		})
	}
}

func NewSearchNode(parent *SearchNode, move Move) *SearchNode {
	return &SearchNode{
		Parent:   parent,
		Children: []*SearchNode{},
		Move:     move,
		Score:    0,
	}
}

// AddChild appends a child node without ordering.
func (n *SearchNode) AddChild(move Move) *SearchNode {
	child := &SearchNode{Parent: n, Move: move, Score: 0, Children: []*SearchNode{}}
	n.Children = append(n.Children, child)
	return child
}

// AddChildScored inserts a child keeping Children ordered (descending) by Score.
// Uses binary search + slice insertion for O(log n + n) complexity; with typical
// chess branching factor (< 40) this is faster than appending then sorting.
func (n *SearchNode) AddChildScored(move Move, score int) *SearchNode {
	child := &SearchNode{Parent: n, Move: move, Score: score, Children: []*SearchNode{}}
	children := n.Children
	// Find first position whose score < new score to maintain descending order.
	idx := sort.Search(len(children), func(i int) bool { return children[i].Score < score })
	// Expand slice and shift
	n.Children = append(children, nil)
	copy(n.Children[idx+1:], n.Children[idx:])
	n.Children[idx] = child
	return child
}

func (n *SearchNode) IsLeaf() bool {
	return len(n.Children) == 0
}

func (n *SearchNode) IsRoot() bool {
	return n.Parent == nil
}

func (n *SearchNode) GetBestMove() *Move {
	var bestChild *SearchNode
	// if n.Board.WhiteToMove {
	bestChild = n.Children[0]
	// } else {
	// 	bestChild = n.Children[len(n.Children)-1]
	// }
	return &bestChild.Move
}

func (n *SearchNode) ToString() string {
	ret := ""
	for _, child := range n.Children {
		ret += fmt.Sprintf("%v Score: %d\n", child.Move.ToString(), child.Score)
	}
	ret += fmt.Sprintf("Total nodes: %d\n", totalNodes)
	return ret
}

func (b *Board) GenerateSearchTree(depth int) *SearchNode {
	totalNodes = 1
	root := NewSearchNode(nil, Move{})
	var generate func(node *SearchNode, isWhite bool, depth int, board *Board)
	generate = func(node *SearchNode, isWhite bool, depth int, board *Board) {
		if depth == 0 {
			return
		}
		moves := board.GetLegalMoves()
		totalNodes += len(moves)
		for _, move := range moves {
			childBoard := board.Clone()
			childBoard.MovePiece(move, board.WhiteToMove)
			if !childBoard.IsKingInCheck(board.WhiteToMove) {
				score := childBoard.Evaluate()
				childNode := node.AddChild(move)
				childNode.Score = score
				generate(childNode, !isWhite, depth-1, childBoard)
			} else {
				// fmt.Printf("Depth %d, Move %s, WTM:%v\n", depth, move.ToString(), childBoard.WhiteToMove)
			}
		}
		if len(node.Children) > 1 {
			// Ordenar hijos para que el mejor esté primero
			node.sortChildren(board.WhiteToMove)
		}
		if len(node.Children) > 0 {
			node.Score = node.Children[0].Score
		}
	}
	generate(root, b.WhiteToMove, depth, b)
	return root
}

func (b *Board) GetBestMove(depth int) (bestMove *Move, bestScore int) {
	root := b.GenerateSearchTree(depth)
	if len(root.Children) == 0 {
		return nil, 0 // No legal moves
	}
	return root.GetBestMove(), root.Children[0].Score
}

func (b *Board) GetBestLine(depth int) (bestMoves MoveList, bestScore int) {
	tree := b.GenerateSearchTree(depth)
	fmt.Println(tree.ToString())
	if len(tree.Children) == 0 {
		return nil, 0 // No legal moves
	}
	bestScore = tree.Children[0].Score
	for tree != nil && len(tree.Children) > 0 {
		bestChild := tree.Children[0]
		move := tree.GetBestMove()
		if move != nil {
			bestMoves.Add(*move)
		}
		bestScore = bestChild.Score
		tree = bestChild
	}
	return bestMoves, bestScore
}

// Equal checks if two boards are identical in piece placement and turn.
func (b *Board) Equal(other *Board) bool {
	return b.WhitePieces == other.WhitePieces &&
		b.BlackPieces == other.BlackPieces &&
		b.WhiteToMove == other.WhiteToMove &&
		b.EnPassant == other.EnPassant &&
		b.Castling == other.Castling
}
