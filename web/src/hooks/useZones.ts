import { useState, useCallback, useEffect } from 'react'
import {
  listZones,
  getZone as getZoneApi,
  createZone as createZoneApi,
  updateZone as updateZoneApi,
  assignPathToZone as assignPathToZoneApi,
} from '../api/client'
import { zoneFromDto, toAssignedAgentsDto } from '../api/mappers'
import type { Zone } from '../api/types'

export interface UseZonesResult {
  zones: Zone[]
  loading: boolean
  error: string | null
  refetch: () => Promise<void>
  getZone: (id: string) => Promise<Zone | null>
  createZone: (params: {
    project_id: string
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

export function useZones(projectId: string | null): UseZonesResult {
  const [zones, setZones] = useState<Zone[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const refetch = useCallback(async () => {
    if (!projectId) {
      setZones([])
      setError(null)
      return
    }
    setLoading(true)
    setError(null)
    try {
      const res = await listZones({ project_id: projectId })
      setZones((res.zones ?? []).map(zoneFromDto).filter(Boolean) as Zone[])
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load zones')
      setZones([])
    } finally {
      setLoading(false)
    }
  }, [projectId])

  const getZone = useCallback(async (id: string): Promise<Zone | null> => {
    try {
      const res = await getZoneApi({ zone_id: id })
      return zoneFromDto(res.zone) ?? null
    } catch {
      return null
    }
  }, [])

  const createZone = useCallback(
    async (params: {
      project_id: string
      name: string
      pattern?: string
      purpose?: string
      constraints?: string[]
      assigned_agent?: string
    }): Promise<Zone | null> => {
      try {
        const res = await createZoneApi({
          project_id: params.project_id,
          name: params.name,
          pattern: params.pattern ?? '',
          purpose: params.purpose ?? '',
          constraints: params.constraints ?? [],
          assigned_agents: toAssignedAgentsDto(params.assigned_agent),
        })
        await refetch()
        return zoneFromDto(res.zone) ?? null
      } catch {
        return null
      }
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
      try {
        const res = await updateZoneApi({
          zone_id: params.zone_id,
          name: params.name ?? '',
          pattern: params.pattern ?? '',
          purpose: params.purpose ?? '',
          constraints: params.constraints ?? [],
          assigned_agents: toAssignedAgentsDto(params.assigned_agent),
        })
        await refetch()
        return zoneFromDto(res.zone) ?? null
      } catch {
        return null
      }
    },
    [refetch]
  )

  const assignPathToZone = useCallback(
    async (zoneId: string, path: string): Promise<Zone | null> => {
      try {
        const res = await assignPathToZoneApi({ zone_id: zoneId, path })
        await refetch()
        return zoneFromDto(res.zone) ?? null
      } catch {
        return null
      }
    },
    [refetch]
  )

  useEffect(() => {
    if (!projectId) {
      setZones([])
      setLoading(false)
      setError(null)
      return
    }
    refetch()
  }, [projectId, refetch])

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
