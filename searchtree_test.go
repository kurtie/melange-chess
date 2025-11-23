package melange

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

func TestSearch(t *testing.T) {
	board := NewBoard()
	root := board.GenerateSearchTree(2)
	assert.Assert(t, root != nil, "Search tree root should not be nil")
	assert.Assert(t, len(root.Children) > 0, "Root should have children")
	assert.Assert(t, !root.IsLeaf(), "Root should not be a leaf")
	assert.Assert(t, root.IsRoot(), "Root should be a root node")

	// Check that children are ordered by score descending
	for i := 1; i < len(root.Children); i++ {
		assert.Assert(t, root.Children[i-1].Score >= root.Children[i].Score, "Children (white) should be ordered by score descending")
	}

	// Check deeper levels
	for _, child := range root.Children {
		assert.Assert(t, len(child.Children) > 0, "Child should have children")
		assert.Assert(t, !child.IsLeaf(), "Child should not be a leaf")
		assert.Assert(t, !child.IsRoot(), "Child should not be a root node")
		for j := 1; j < len(child.Children); j++ {
			assert.Assert(t, child.Children[j-1].Score <= child.Children[j].Score, "Grandchildren (black) should be ordered by score ascending")
		}
	}
}

func TestSearchWithHash(t *testing.T) {
	board := NewBoard()
	// board.SetFen("7k/8/8/8/8/7P/8/7K w - - 0 1")
	fmt.Println("Initial Board:", board.ToString())
	root := board.GenerateSearchTree(4)
	assert.Assert(t, root != nil, "Search tree root should not be nil")
	fmt.Println("Search Tree:\n", root.ToString())
	assert.Assert(t, false == true, "Force output to be seen")
}
