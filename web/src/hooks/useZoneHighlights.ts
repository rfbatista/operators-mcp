import { useState, useEffect, useCallback } from 'react'
import { callTool } from '../api/callTool'
import type { Zone } from '../api/types'

export interface ZoneHighlights {
  highlightPaths: Set<string>
  pathToZones: Map<string, string[]>
  zones: Zone[]
  loading: boolean
  error: string | null
}

export function useZoneHighlights(): ZoneHighlights & { refetch: () => void } {
  const [zones, setZones] = useState<Zone[]>([])
  const [highlightPaths, setHighlightPaths] = useState<Set<string>>(new Set())
  const [pathToZones, setPathToZones] = useState<Map<string, string[]>>(new Map())
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchHighlights = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const listRes = await callTool('list_zones', {})
      if (listRes.isError || !listRes.content?.[0]?.text) {
        setZones([])
        setHighlightPaths(new Set())
        setPathToZones(new Map())
        return
      }
      const { zones: zs } = JSON.parse(listRes.content[0].text) as { zones: Zone[] }
      setZones(zs)

      const allPaths = new Set<string>()
      const pathToZonesMap = new Map<string, string[]>()

      for (const zone of zs) {
        if (!zone.pattern?.trim()) {
          for (const p of zone.explicit_paths || []) {
            allPaths.add(p)
            pathToZonesMap.set(p, [...(pathToZonesMap.get(p) ?? []), zone.name])
          }
          continue
        }
        const matchRes = await callTool('list_matching_paths', {
          pattern: zone.pattern,
        })
        if (matchRes.isError || !matchRes.content?.[0]?.text) continue
        const { paths } = JSON.parse(matchRes.content[0].text) as { paths: string[] }
        for (const p of paths) {
          allPaths.add(p)
          pathToZonesMap.set(p, [...(pathToZonesMap.get(p) ?? []), zone.name])
        }
        for (const p of zone.explicit_paths || []) {
          allPaths.add(p)
          pathToZonesMap.set(p, [...(pathToZonesMap.get(p) ?? []), zone.name])
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
  }, [])

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
