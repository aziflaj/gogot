package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aziflaj/gogot/fileutils"
)

// IndexTree ...
type IndexTree struct {
	Name     string
	Hash     string
	Children []*IndexTree
}

func BuildIndexFromFile(file *os.File) *IndexTree {
	tree := &IndexTree{}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		splitLine := strings.Split(scanner.Text(), " ")
		hash := splitLine[0]
		path := splitLine[1]

		tree.AddPath(path, hash)
	}

	return tree
}

func BuildIndexFromCommit(hash string, name string) (tree *IndexTree, err error) {
	content, err := fileutils.ReadCommitContents(hash)
	if err != nil {
		return
	}

	tree = &IndexTree{Name: name}

	for _, line := range strings.Split(content, "\n") {
		splitContent := strings.Split(line, " ")
		if len(splitContent) == 1 { // newline at the end of file
			continue
		}

		objectType := splitContent[0]
		hash := splitContent[1]
		objectName := splitContent[2]

		if objectType == "blob" { // Object is a file
			tree.Children = append(tree.Children, &IndexTree{Name: objectName, Hash: hash})
		} else if objectType == "tree" { // Object is a dir
			child, _ := BuildIndexFromCommit(hash, objectName)
			tree.Children = append(tree.Children, child)
		}
	}

	return
}

func NewTreeWithName(name string) *IndexTree {
	return &IndexTree{Name: name}
}

func (t *IndexTree) FindChildByName(name string) *IndexTree {
	if t.Name == name {
		return t
	}

	if len(t.Children) == 0 {
		return nil
	}

	// depth-first search for child
	for _, child := range t.Children {
		if child.Name == name {
			return child
		}

		if namedChild := child.FindChildByName(name); namedChild != nil {
			return namedChild
		}
	}

	return nil
}

func (t *IndexTree) FindChildByPath(path string) *IndexTree {
	pathParts := strings.Split(path, "/")
	child := t.FindChildByName(pathParts[0])
	if child == nil {
		return nil
	}

	if len(pathParts) == 1 {
		return child
	}

	return child.FindChildByPath(strings.Join(pathParts[1:], "/"))
}

// AddPath ...
func (t *IndexTree) AddPath(path string, hash string) {
	splitPath := strings.Split(path, "/")
	for idx, pathPart := range splitPath {
		// fmt.Println("pathPart: " + pathPart)

		if t.Name == "" {
			// fmt.Println("Root with no name, adding " + pathPart + " as root")
			t.Name = pathPart
		} else if child := t.FindChildByName(pathPart); child != nil {
			restOfPath := strings.Split(path, pathPart+"/")
			// fmt.Println("Found child " + pathPart + ", adding " + restOfPath[1] + " to " + child.Name)
			// fmt.Println(path, pathPart, restOfPath[1])
			child.AddPath(restOfPath[1], hash)
			return
		} else if idx < len(splitPath)-1 {
			// fmt.Println("Appending non-leaf child " + pathPart)
			child := NewTreeWithName(pathPart)
			child.AddPath(path, hash)
			t.Children = append(t.Children, child)
			return
		} else {
			// fmt.Println("Appending last child (leaf) " + pathPart + " in " + t.Name)
			t.Children = append(t.Children, &IndexTree{Name: pathPart, Hash: hash})
		}
	}

	fmt.Printf("\n\n\n")
}

func (t *IndexTree) BuildObjectTree(name string) (string, error) {
	hash := TimedHash(name)
	file, err := fileutils.CreateAndOpenCommitFile(hash)
	if err != nil {
		return "", err
	}

	defer file.Close()

	for _, child := range t.Children {
		if child.Hash != "" { // is a file
			file.WriteString(fmt.Sprintf("blob %s %s\n", child.Hash, child.Name))
		} else { // is a dir
			dirHash, err := child.BuildObjectTree(child.Name)
			if err != nil {
				return "", err
			}
			file.WriteString(fmt.Sprintf("tree %s %s\n", dirHash, child.Name))
		}
	}

	return hash, nil
}

func (t IndexTree) String() string {
	if t.Name == "" {
		return "[EMPTY]"
	}

	var builder strings.Builder

	if len(t.Children) == 0 {
		builder.WriteString(fmt.Sprintf("* '%s': no children\n", t.Name))
		return builder.String()
	}

	var childrenNames []string
	builder.WriteString(fmt.Sprintf("* %s: \n", t.Name))
	for _, child := range t.Children {
		childrenNames = append(childrenNames, child.Name)
		builder.WriteString(fmt.Sprintf("    -> %s (%s): \n", child.Name, child.Hash))
	}

	for _, child := range t.Children {
		builder.WriteString(child.String())
	}

	return builder.String()
}
