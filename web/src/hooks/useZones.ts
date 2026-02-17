import { useState, useCallback, useEffect } from 'react'
import { callTool } from '../api/callTool'
import type { Zone } from '../api/types'

export interface UseZonesResult {
  zones: Zone[]
  loading: boolean
  error: string | null
  refetch: () => Promise<void>
  getZone: (id: string) => Promise<Zone | null>
  createZone: (params: {
    name: string
    pattern?: string
    purpose?: string
    constraints?: string[]
    assigned_agent?: string
  }) => Promise<Zone | null>
  updateZone: (params: {
    zone_id: string
    name?: string
    pattern?: string
    purpose?: string
    constraints?: string[]
    assigned_agent?: string
  }) => Promise<Zone | null>
  assignPathToZone: (zoneId: string, path: string) => Promise<Zone | null>
}

export function useZones(): UseZonesResult {
  const [zones, setZones] = useState<Zone[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refetch = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await callTool('list_zones', {})
      if (res.isError || !res.content?.[0]?.text) {
        setZones([])
        return
      }
      const data = JSON.parse(res.content[0].text) as { zones: Zone[] }
      setZones(data.zones ?? [])
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load zones')
      setZones([])
    } finally {
      setLoading(false)
    }
  }, [])

  const getZone = useCallback(async (id: string): Promise<Zone | null> => {
    const res = await callTool('get_zone', { zone_id: id })
    if (res.isError || !res.content?.[0]?.text) return null
    const data = JSON.parse(res.content[0].text) as { zone: Zone }
    return data.zone ?? null
  }, [])

  const createZone = useCallback(
    async (params: {
      name: string
      pattern?: string
      purpose?: string
      constraints?: string[]
      assigned_agent?: string
    }): Promise<Zone | null> => {
      const res = await callTool('create_zone', {
        name: params.name,
        pattern: params.pattern ?? '',
        purpose: params.purpose ?? '',
        constraints: params.constraints ?? [],
        assigned_agent: params.assigned_agent ?? '',
      })
      if (res.isError || !res.content?.[0]?.text) return null
      const data = JSON.parse(res.content[0].text) as { zone: Zone }
      await refetch()
      return data.zone ?? null
    },
    [refetch]
  )

  const updateZone = useCallback(
    async (params: {
      zone_id: string
      name?: string
      pattern?: string
      purpose?: string
      constraints?: string[]
      assigned_agent?: string
    }): Promise<Zone | null> => {
      const res = await callTool('update_zone', {
        zone_id: params.zone_id,
        name: params.name ?? '',
        pattern: params.pattern ?? '',
        purpose: params.purpose ?? '',
        constraints: params.constraints ?? [],
        assigned_agent: params.assigned_agent ?? '',
      })
      if (res.isError || !res.content?.[0]?.text) return null
      const data = JSON.parse(res.content[0].text) as { zone: Zone }
      await refetch()
      return data.zone ?? null
    },
    [refetch]
  )

  const assignPathToZone = useCallback(
    async (zoneId: string, path: string): Promise<Zone | null> => {
      const res = await callTool('assign_path_to_zone', { zone_id: zoneId, path })
      if (res.isError || !res.content?.[0]?.text) return null
      const data = JSON.parse(res.content[0].text) as { zone: Zone }
      await refetch()
      return data.zone ?? null
    },
    [refetch]
  )

  useEffect(() => {
    refetch()
  }, [refetch])

  return {
    zones,
    loading,
    error,
    refetch,
    getZone,
    createZone,
    updateZone,
    assignPathToZone,
  }
}
