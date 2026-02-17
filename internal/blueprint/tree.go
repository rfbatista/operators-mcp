package blueprint

import (
	"os"
	"path/filepath"
)

// TreeNode is a node in the source tree (path, name, is_dir, children).
type TreeNode struct {
	Path     string      `json:"path"`
	Name     string      `json:"name"`
	IsDir    bool        `json:"is_dir"`
	Children []*TreeNode `json:"children"`
}

// ListTree builds a tree from root. Root empty means cwd. Returns error if root unreadable.
func ListTree(root string) (*TreeNode, error) {
	if root == "" {
		var err error
		root, err = os.Getwd()
		if err != nil {
			return nil, &StructuredError{Code: "ROOT_UNREADABLE", Message: err.Error()}
		}
	}
	info, err := os.Stat(root)
	if err != nil {
		return nil, &StructuredError{Code: "ROOT_UNREADABLE", Message: err.Error()}
	}
	if !info.IsDir() {
		return nil, &StructuredError{Code: "ROOT_UNREADABLE", Message: "root is not a directory"}
	}
	return listTreeAt(root, root, "")
}

func listTreeAt(fullRoot, current, relPath string) (*TreeNode, error) {
	entries, err := os.ReadDir(current)
	if err != nil {
		return nil, &StructuredError{Code: "ROOT_UNREADABLE", Message: err.Error()}
	}
	name := filepath.Base(current)
	if relPath == "" {
		name = "."
	}
	node := &TreeNode{
		Path:     filepath.ToSlash(relPath),
		Name:     name,
		IsDir:    true,
		Children: nil,
	}
	for _, e := range entries {
		childRel := filepath.Join(relPath, e.Name())
		childRel = filepath.ToSlash(childRel)
		childFull := filepath.Join(current, e.Name())
		if e.IsDir() {
			child, err := listTreeAt(fullRoot, childFull, childRel)
			if err != nil {
				return nil, err
			}
			node.Children = append(node.Children, child)
		} else {
			node.Children = append(node.Children, &TreeNode{
				Path:     childRel,
				Name:     e.Name(),
				IsDir:    false,
				Children: nil,
			})
		}
	}
	return node, nil
}
