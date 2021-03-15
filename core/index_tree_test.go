package core

import (
	"testing"
)

func generateTree() IndexTree {
	return IndexTree{
		Name: ".",
		Children: []*IndexTree{
			{Name: "commands",
				Children: []*IndexTree{
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

func TestFindChildByName(t *testing.T) {
	tree := generateTree()
	commandsChild := tree.FindChildByName("commands")
	if commandsChild == nil {
		t.Error("Expected ChildWithName to return commands, returned nil")
	}

	root := tree.FindChildByName(".")
	if root == nil {
		t.Error("Expected ChildWithName to return root, returned nil")
	}
}

func TestAddRootFilePath(t *testing.T) {
	tree := generateTree()
	tree.AddPath("./go.mod", "hash")

	modFile := tree.FindChildByName("go.mod")
	if modFile == nil {
		t.Error("Expected modFile to not be nil")
	}
}

func TestAddNestedFilePath(t *testing.T) {
	tree := generateTree()
	tree.AddPath("./indextree/tree.go", "hash")
	indextree := tree.FindChildByName("indextree")
	if indextree == nil {
		t.Error("Expected indextree to not be nil")
	}

	treeFile := indextree.FindChildByName("tree.go")
	if treeFile == nil {
		t.Error("Expected treeFile to not be nil")
	}
}

func TestAddFileInExistingPath(t *testing.T) {
	tree := generateTree()
	tree.AddPath("./commands/remote.go", "hash")

	remoteFile := tree.FindChildByName("remote.go")
	if remoteFile == nil {
		t.Error("Expected remoteFile to not be nil")
	}
}
