package indextree

import (
	"fmt"
	"testing"
)

func generateTree() IndexTree {
	return IndexTree{
		Name: ".",
		Children: []IndexTree{
			{Name: "commands",
				Children: []IndexTree{
					{Name: "add.go"},
					{Name: "commit.go"},
					{Name: "init.go"},
				},
			},
			{Name: "gogot.go"},
			{Name: ".gitignore"},
		},
	}
}

func TestChildWithName(t *testing.T) {
	tree := generateTree()
	commandsChild := tree.ChildWithName("commands")
	if commandsChild == nil {
		t.Error("Expected ChildWithName to return commands, returned nil")
	}

	root := tree.ChildWithName(".")
	if root == nil {
		t.Error("Expected ChildWithName to return root, returned nil")
	}
}

func TestAddRootFilePath(t *testing.T) {
	tree := generateTree()
	tree.AddPath("./go.mod", "sha")
	// fmt.Println(tree)

	modFile := tree.ChildWithName("go.mod")
	if modFile == nil {
		t.Error("Expected modFile to not be nil")
	}

	fmt.Println(tree)
}

func TestAddNestedFilePath(t *testing.T) {
	tree := generateTree()
	tree.AddPath("./indextree/tree.go", "sha")
	// fmt.Println(tree)

	indextree := tree.ChildWithName("indextree")
	if indextree == nil {
		t.Error("Expected indextree to not be nil")
	}

	treeFile := indextree.ChildWithName("tree.go")
	if treeFile == nil {
		t.Error("Expected tree to not be nil")
	}
}
