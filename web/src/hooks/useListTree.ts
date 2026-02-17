import { useState, useEffect, useCallback } from 'react'
import { listTree } from '../api/client'
import { treeNodeFromDto } from '../api/mappers'
import type { TreeNode } from '../api/types'

export interface ListTreeState {
  tree: TreeNode | null
  loading: boolean
  error: string | null
}

export function useListTree(root?: string): ListTreeState & { refetch: () => void } {
  const [tree, setTree] = useState<TreeNode | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchTree = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await listTree(root != null ? { root } : {})
      setTree(treeNodeFromDto(res.tree) ?? null)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load tree')
      setTree(null)
    } finally {
      setLoading(false)
    }
  }, [root])

  useEffect(() => {
    fetchTree()
  }, [fetchTree])

  return { tree, loading, error, refetch: fetchTree }
}
