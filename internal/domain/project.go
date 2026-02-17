package domain

// Project defines the directory root that everything (tree, matching paths, zones) is based on.
// All paths and operations are relative to the project's root.
// IgnoredPaths are paths (files or directories) to hide from the tree view; children of ignored dirs are hidden too.
type Project struct {
	ID           string
	Name         string
	RootDir      string
	IgnoredPaths []string
}
