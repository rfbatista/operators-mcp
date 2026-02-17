import { useState, useEffect, useCallback } from 'react'
import { listZones, listMatchingPaths, listTree } from '../api/client'
import { treeNodeFromDto, zoneFromDto } from '../api/mappers'
import type { Zone } from '../api/types'
import type { TreeNode } from '../api/types'

export interface ZoneHighlights {
  highlightPaths: Set<string>
  pathToZones: Map<string, string[]>
  zones: Zone[]
  loading: boolean
  error: string | null
}

/** Flatten tree to all paths (node path + all descendants). */
function flattenTreePaths(node: TreeNode): string[] {
  const paths = [node.path]
  for (const child of node.children ?? []) {
    paths.push(...flattenTreePaths(child))
  }
  return paths
}

/** Path is the assigned path or a descendant of it (direct child or nested). */
function pathIsAssignedOrDescendant(path: string, assignedPath: string): boolean {
  if (path === assignedPath) return true
  const prefix = assignedPath.endsWith('/') ? assignedPath : assignedPath + '/'
  return path.startsWith(prefix)
}

export function useZoneHighlights(projectId: string | null): ZoneHighlights & { refetch: () => void } {
  const [zones, setZones] = useState<Zone[]>([])
  const [highlightPaths, setHighlightPaths] = useState<Set<string>>(new Set())
  const [pathToZones, setPathToZones] = useState<Map<string, string[]>>(new Map())
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetchHighlights = useCallback(async () => {
    if (!projectId) {
      setZones([])
      setHighlightPaths(new Set())
      setPathToZones(new Map())
      setError(null)
      return
    }
    setLoading(true)
    setError(null)
    try {
      const [listRes, treeRes] = await Promise.all([
        listZones({ project_id: projectId }),
        listTree({ project_id: projectId }),
      ])
      const zs = (listRes.zones ?? []).map(zoneFromDto).filter(Boolean) as Zone[]
      setZones(zs)

      const tree = treeNodeFromDto(treeRes.tree)
      const allTreePaths = tree ? flattenTreePaths(tree) : []

      const allPaths = new Set<string>()
      const pathToZonesMap = new Map<string, string[]>()

      for (const zone of zs) {
        // Explicit paths: the assigned path and all its children (directories and files)
        for (const assignedPath of zone.explicit_paths ?? []) {
          const normalized = assignedPath.trim()
          if (!normalized) continue
          for (const p of allTreePaths) {
            if (pathIsAssignedOrDescendant(p, normalized)) {
              allPaths.add(p)
              pathToZonesMap.set(p, [...(pathToZonesMap.get(p) ?? []), zone.name])
            }
          }
          // Ensure the assigned path itself is in the set even if not in tree yet
          if (!allPaths.has(normalized)) {
            allPaths.add(normalized)
            pathToZonesMap.set(normalized, [...(pathToZonesMap.get(normalized) ?? []), zone.name])
          }
        }

        if (!zone.pattern?.trim()) continue

        try {
          const matchRes = await listMatchingPaths({
            pattern: zone.pattern,
            project_id: projectId,
          })
          const paths = matchRes.paths ?? []
          for (const p of paths) {
            allPaths.add(p)
            pathToZonesMap.set(p, [...(pathToZonesMap.get(p) ?? []), zone.name])
          }
        } catch {
          // skip zone on match error
        }
      }

      setHighlightPaths(allPaths)
      setPathToZones(pathToZonesMap)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load zone highlights')
      setHighlightPaths(new Set())
      setPathToZones(new Map())
    } finally {
      setLoading(false)
    }
  }, [projectId])

  useEffect(() => {
    fetchHighlights()
  }, [fetchHighlights])

  return {
    highlightPaths,
    pathToZones,
    zones,
    loading,
    error,
    refetch: fetchHighlights,
  }
}
