import { useState, useEffect, useCallback } from 'react'
import { listZones, listMatchingPaths } from '../api/client'
import { zoneFromDto } from '../api/mappers'
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
      const listRes = await listZones()
      const zs = (listRes.zones ?? []).map(zoneFromDto).filter(Boolean) as Zone[]
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
        try {
          const matchRes = await listMatchingPaths({ pattern: zone.pattern })
          const paths = matchRes.paths ?? []
          for (const p of paths) {
            allPaths.add(p)
            pathToZonesMap.set(p, [...(pathToZonesMap.get(p) ?? []), zone.name])
          }
        } catch {
          // skip zone on match error
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
