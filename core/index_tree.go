package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

	// fmt.Printf("\n\n\n")
}

func (t *IndexTree) BuildObjectTree(name string) (string, error) {
	objectTree, err := CreateObjectFromString(name)
	if err != nil {
		return "", err
	}

	defer objectTree.FlushAndClose()

	for _, child := range t.Children {
		if child.Hash != "" { // is a file
			objectTree.AddBlob(child)
		} else { // is a dir
			dirHash, err := child.BuildObjectTree(child.Name)
			if err != nil {
				return "", err
			}
			objectTree.AddTree(child, dirHash)
		}
	}

	return objectTree.Hash, nil
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