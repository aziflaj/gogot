package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aziflaj/gogot/indextree"
)

// Commit ...
func Commit(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: gogot commit -m \"message\"")
		os.Exit(1)
	}

	indexTree := buildIndexTree()
	buildObjectTree("root", *indexTree)

	// root_sha = build_tree("root", index_tree)
	// commit_sha = build_commit(tree: root_sha)
	// update_ref(commit_sha: commit_sha)
	// clear_index
}

// func buildTree(name string, tree Tree) {

// }

func buildIndexTree() *indextree.IndexTree {
	indexTree := indextree.New()

	indexFile, err := os.Open(indexPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer indexFile.Close()
	scanner := bufio.NewScanner(indexFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		splitLine := strings.Split(scanner.Text(), " ")
		sha := splitLine[0]
		path := splitLine[1]
		indexTree.AddPath(path, sha)
	}

	// fmt.Println(indexTree)
	return indexTree
}

func buildObjectTree(name string, tree indextree.IndexTree) {
	fmt.Println(tree)
}

// gogotTree := NewTree()
// gogotTree.AddPath("./path/to/file1.txt", "sha1")
// gogotTree.AddPath("./path/at/file2.txt", "sha1")
// gogotTree.AddPath("./path/to/file/3.txt", "sha1")
