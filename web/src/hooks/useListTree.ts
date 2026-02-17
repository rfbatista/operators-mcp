import { useState, useEffect, useCallback } from 'react'
import { callTool } from '../api/callTool'
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
      const res = await callTool('list_tree', root != null ? { root } : {})
      if (res.isError || !res.content?.[0]?.text) {
        setError('Failed to load tree')
        setTree(null)
        return
      }
      const data = JSON.parse(res.content[0].text) as { tree: TreeNode }
      setTree(data.tree)
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
