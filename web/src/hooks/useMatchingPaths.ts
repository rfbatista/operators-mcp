import { useState, useEffect, useCallback, useRef } from 'react'
import { callTool } from '../api/callTool'

const DEBOUNCE_MS = 400

export interface UseMatchingPathsResult {
  pattern: string
  setPattern: (p: string) => void
  paths: string[]
  loading: boolean
  error: string | null
  invalidPattern: boolean
}

export function useMatchingPaths(): UseMatchingPathsResult {
  const [pattern, setPattern] = useState('')
  const [paths, setPaths] = useState<string[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [invalidPattern, setInvalidPattern] = useState(false)
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  const fetchPaths = useCallback(async (p: string) => {
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
      const res = await callTool('list_matching_paths', { pattern: p })
      if (res.isError) {
        const text = res.content?.[0]?.text ?? ''
        if (text.includes('INVALID_PATTERN') || /invalid|regex|syntax/i.test(text)) {
          setInvalidPattern(true)
          setPaths([])
        } else {
          setError('Failed to get matching paths')
        }
        return
      }
      if (!res.content?.[0]?.text) {
        setPaths([])
        return
      }
      const data = JSON.parse(res.content[0].text) as { paths?: string[] }
      setPaths(data.paths ?? [])
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to get matching paths')
      setPaths([])
    } finally {
      setLoading(false)
    }
  }, [])

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
      fetchPaths(pattern)
    }, DEBOUNCE_MS)
    return () => {
      if (timerRef.current) clearTimeout(timerRef.current)
    }
  }, [pattern, fetchPaths])

  return {
    pattern,
    setPattern,
    paths,
    loading,
    error,
    invalidPattern,
  }
}
