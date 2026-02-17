import { useState, useEffect, useCallback } from 'react'
import { listTree } from '../api/client'
import { treeNodeFromDto } from '../api/mappers'
import type { TreeNode } from '../api/types'

export interface ListTreeState {
  tree: TreeNode | null
  loading: boolean
  error: string | null
}

export function useListTree(projectId: string | null, root?: string): ListTreeState & { refetch: () => void } {
  const [tree, setTree] = useState<TreeNode | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetchTree = useCallback(async () => {
    if (!projectId && (root == null || root === '')) {
      setTree(null)
      setLoading(false)
      setError(null)
      return
    }
    setLoading(true)
    setError(null)
    try {
      const res = await listTree(
        projectId ? { project_id: projectId } : root != null ? { root } : {}
      )
      setTree(treeNodeFromDto(res.tree) ?? null)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load tree')
      setTree(null)
    } finally {
      setLoading(false)
    }
  }, [projectId, root])

  useEffect(() => {
    fetchTree()
  }, [fetchTree])

  return { tree, loading, error, refetch: fetchTree }
}
