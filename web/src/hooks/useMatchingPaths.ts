import { useState, useEffect, useCallback, useRef } from 'react'
import { listMatchingPaths } from '../api/client'
import { ApiError } from '../api/client'

const DEBOUNCE_MS = 400

export interface UseMatchingPathsResult {
  pattern: string
  setPattern: (p: string) => void
  paths: string[]
  loading: boolean
  error: string | null
  invalidPattern: boolean
}

export function useMatchingPaths(projectId: string | null): UseMatchingPathsResult {
  const [pattern, setPattern] = useState('')
  const [paths, setPaths] = useState<string[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [invalidPattern, setInvalidPattern] = useState(false)
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  const fetchPaths = useCallback(async (p: string, projectId: string | null) => {
    if (!p.trim()) {
      setPaths([])
      setInvalidPattern(false)
      setError(null)
      return
    }
    setLoading(true)
    setError(null)
    setInvalidPattern(false)
    try {
      const res = await listMatchingPaths({
        pattern: p,
        ...(projectId ? { project_id: projectId } : {}),
      })
      setPaths(res.paths ?? [])
    } catch (e) {
      const msg = e instanceof Error ? e.message : ''
      if (
        e instanceof ApiError &&
        (e.message.includes('INVALID_PATTERN') || /invalid|regex|syntax/i.test(e.message))
      ) {
        setInvalidPattern(true)
        setPaths([])
      } else {
        setError(msg || 'Failed to get matching paths')
      }
    } finally {
      setLoading(false)
    }
  }, [projectId])

  useEffect(() => {
    if (timerRef.current) clearTimeout(timerRef.current)
    if (!pattern.trim()) {
      setPaths([])
      setLoading(false)
      setInvalidPattern(false)
      setError(null)
      return
    }
    timerRef.current = setTimeout(() => {
      timerRef.current = null
      fetchPaths(pattern, projectId)
    }, DEBOUNCE_MS)
    return () => {
      if (timerRef.current) clearTimeout(timerRef.current)
    }
  }, [pattern, projectId, fetchPaths])

  return {
    pattern,
    setPattern,
    paths,
    loading,
    error,
    invalidPattern,
  }
}
