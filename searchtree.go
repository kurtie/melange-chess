package muabdib

type SearchTree struct {
	Root *SearchNode
}

type SearchNode struct {
	Parent   *SearchNode
	Children []*SearchNode
	Board    *Board
	Score    float32
}
