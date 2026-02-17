import { useState, useMemo } from 'react'
import type { TreeNode } from '../api/types'

export interface SourceTreeProps {
  tree: TreeNode | null
  /** Set of relative paths that should be highlighted (e.g. matching a zone) */
  highlightPaths?: Set<string>
  /** Optional: path -> zone name(s) for tooltip or multi-zone highlight */
  pathToZones?: Map<string, string[]>
  onAssignToZone?: (path: string) => void
}

function TreeNodeRow({
  node,
  depth,
  highlightPaths,
  pathToZones,
  onAssignToZone,
}: {
  node: TreeNode
  depth: number
  highlightPaths?: Set<string>
  pathToZones?: Map<string, string[]>
  onAssignToZone?: (path: string) => void
}) {
  const [expanded, setExpanded] = useState(depth < 2)
  const hasChildren = node.children && node.children.length > 0
  const isHighlighted = highlightPaths?.has(node.path) ?? false
  const zones = pathToZones?.get(node.path)

  return (
    <div style={{ marginLeft: depth * 12 }}>
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          gap: 4,
          padding: '2px 4px',
          borderRadius: 4,
          backgroundColor: isHighlighted ? 'rgba(100, 149, 237, 0.25)' : undefined,
        }}
      >
        {hasChildren ? (
          <button
            type="button"
            onClick={() => setExpanded((e) => !e)}
            style={{ padding: 0, marginRight: 4, cursor: 'pointer' }}
            aria-label={expanded ? 'Collapse' : 'Expand'}
          >
            {expanded ? '▼' : '▶'}
          </button>
        ) : (
          <span style={{ width: 14, display: 'inline-block' }} />
        )}
        <span title={zones?.length ? `Zones: ${zones.join(', ')}` : undefined}>
          {node.name}
        </span>
        {node.is_dir && onAssignToZone && (
          <button
            type="button"
            onClick={() => onAssignToZone(node.path)}
            style={{ marginLeft: 'auto', fontSize: 11 }}
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
  onAssignToZone,
}: SourceTreeProps) {
  const safePaths = useMemo(
    () => (highlightPaths ? new Set(highlightPaths) : undefined),
    [highlightPaths]
  )

  if (!tree) return null

  return (
    <div style={{ fontFamily: 'monospace', fontSize: 13 }}>
      <TreeNodeRow
        node={tree}
        depth={0}
        highlightPaths={safePaths}
        pathToZones={pathToZones}
        onAssignToZone={onAssignToZone}
      />
    </div>
  )
}
