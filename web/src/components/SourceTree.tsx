import { useState, useMemo } from 'react'
import type { TreeNode } from '../api/types'

/** Returns true if the node path should be hidden (is ignored or under an ignored dir). */
function isPathIgnored(nodePath: string, ignoredPaths: Set<string>): boolean {
  if (ignoredPaths.size === 0) return false
  for (const ig of ignoredPaths) {
    if (nodePath === ig || nodePath.startsWith(ig + '/')) return true
  }
  return false
}

/** Filter tree recursively: remove nodes that are ignored or under an ignored path. */
function filterTreeByIgnored(node: TreeNode, ignoredPaths: Set<string>): TreeNode | null {
  if (isPathIgnored(node.path, ignoredPaths)) return null
  if (!node.children?.length) return node
  const filteredChildren = node.children
    .map((c) => filterTreeByIgnored(c, ignoredPaths))
    .filter((c): c is TreeNode => c != null)
  return { ...node, children: filteredChildren }
}

export interface SourceTreeProps {
  tree: TreeNode | null
  /** Set of relative paths that should be highlighted (e.g. matching a zone) */
  highlightPaths?: Set<string>
  /** Optional: path -> zone name(s) for tooltip or multi-zone highlight */
  pathToZones?: Map<string, string[]>
  /** Paths to hide from the tree (e.g. project ignored_paths) */
  ignoredPaths?: string[]
  onAssignToZone?: (path: string) => void
  /** Called when user chooses to hide a file/directory */
  onHide?: (path: string) => void
  /** Called when user chooses to show a previously ignored path (for use in "Ignored" list) */
  onShow?: (path: string) => void
}

function TreeNodeRow({
  node,
  depth,
  highlightPaths,
  pathToZones,
  onAssignToZone,
  onHide,
}: {
  node: TreeNode
  depth: number
  highlightPaths?: Set<string>
  pathToZones?: Map<string, string[]>
  onAssignToZone?: (path: string) => void
  onHide?: (path: string) => void
}) {
  const [expanded, setExpanded] = useState(depth < 2)
  const hasChildren = node.children && node.children.length > 0
  const isHighlighted = highlightPaths?.has(node.path) ?? false
  const zones = pathToZones?.get(node.path)

  return (
    <div style={{ marginLeft: depth * 12 }}>
      <div
        className={`flex items-center gap-1 rounded px-1 py-0.5 transition-colors ${
          isHighlighted ? 'bg-primary/20' : ''
        } hover:bg-base-300/70`}
      >
        {hasChildren ? (
          <button
            type="button"
            className="btn btn-ghost btn-xs min-h-6 min-w-6 p-0"
            onClick={() => setExpanded((e) => !e)}
            aria-label={expanded ? 'Collapse' : 'Expand'}
          >
            {expanded ? '▼' : '▶'}
          </button>
        ) : (
          <span className="inline-block w-3.5" />
        )}
        <span
          className="flex-1 font-mono text-sm"
          title={zones?.length ? `Zones: ${zones.join(', ')}` : undefined}
        >
          {node.name}
        </span>
        {onHide && (
          <button
            type="button"
            className="btn btn-ghost btn-xs text-xs opacity-70 hover:opacity-100"
            onClick={() => onHide(node.path)}
            title="Hide from tree"
          >
            Hide
          </button>
        )}
        {node.is_dir && onAssignToZone && (
          <button
            type="button"
            className="btn btn-ghost btn-xs text-xs"
            onClick={() => onAssignToZone(node.path)}
          >
            Assign to zone
          </button>
        )}
      </div>
      {expanded && hasChildren && node.children && (
        <div>
          {node.children.map((child) => (
            <TreeNodeRow
              key={child.path || '.'}
              node={child}
              depth={depth + 1}
              highlightPaths={highlightPaths}
              pathToZones={pathToZones}
              onAssignToZone={onAssignToZone}
              onHide={onHide}
            />
          ))}
        </div>
      )}
    </div>
  )
}

export function SourceTree({
  tree,
  highlightPaths,
  pathToZones,
  ignoredPaths = [],
  onAssignToZone,
  onHide,
  onShow,
}: SourceTreeProps) {
  const safePaths = useMemo(
    () => (highlightPaths ? new Set(highlightPaths) : undefined),
    [highlightPaths]
  )
  const ignoredSet = useMemo(
    () => new Set(ignoredPaths),
    [ignoredPaths]
  )
  const filteredTree = useMemo(() => {
    if (!tree) return null
    if (ignoredSet.size === 0) return tree
    return filterTreeByIgnored(tree, ignoredSet)
  }, [tree, ignoredSet])

  if (!filteredTree) return null

  return (
    <div className="font-mono text-sm">
      {ignoredPaths.length > 0 && onShow && (
        <div className="mb-2 rounded bg-base-200/60 px-2 py-1.5 text-xs">
          <span className="text-base-content/70">Hidden: </span>
          {ignoredPaths.map((path) => (
            <span key={path} className="mr-1 inline-flex items-center gap-0.5">
              <code className="rounded bg-base-300 px-1">{path}</code>
              <button
                type="button"
                className="btn btn-ghost btn-xs min-h-5 px-1 text-xs"
                onClick={() => onShow(path)}
              >
                Show
              </button>
            </span>
          ))}
        </div>
      )}
      <TreeNodeRow
        node={filteredTree}
        depth={0}
        highlightPaths={safePaths}
        pathToZones={pathToZones}
        onAssignToZone={onAssignToZone}
        onHide={onHide}
      />
    </div>
  )
}
