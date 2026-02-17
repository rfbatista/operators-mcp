package domain

// TreeNode is a node in the source tree (path, name, is_dir, children).
// Used when listing project structure for the designer.
type TreeNode struct {
	Path     string
	Name     string
	IsDir    bool
	Children []*TreeNode
}
