import { useState, useCallback, useEffect } from 'react'
import { listAgents, createAgent, updateAgent, deleteAgent } from '../api/client'
import { agentFromDto } from '../api/mappers'
import type { Agent } from '../api/types'

export interface UseAgentsResult {
  agents: Agent[]
  loading: boolean
  error: string | null
  refetch: () => Promise<void>
  createAgent: (params: { name?: string; description?: string; prompt?: string }) => Promise<Agent | null>
  updateAgent: (agentId: string, params: { name?: string; description?: string; prompt?: string }) => Promise<Agent | null>
  deleteAgent: (agentId: string) => Promise<void>
}

export function useAgents(): UseAgentsResult {
  const [agents, setAgents] = useState<Agent[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refetch = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await listAgents()
      setAgents((res.agents ?? []).map(agentFromDto).filter(Boolean) as Agent[])
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load agents')
      setAgents([])
    } finally {
      setLoading(false)
    }
  }, [])

  const createAgentApi = useCallback(
    async (params: { name?: string; description?: string; prompt?: string }): Promise<Agent | null> => {
      try {
        const res = await createAgent({
          name: params.name ?? '',
          description: params.description ?? '',
          prompt: params.prompt ?? '',
        })
        await refetch()
        return agentFromDto(res.agent) ?? null
      } catch {
        return null
      }
    },
    [refetch]
  )

  const updateAgentApi = useCallback(
    async (agentId: string, params: { name?: string; description?: string; prompt?: string }): Promise<Agent | null> => {
      try {
        const res = await updateAgent({
          agent_id: agentId,
          name: params.name ?? '',
          description: params.description ?? '',
          prompt: params.prompt ?? '',
        })
        await refetch()
        return agentFromDto(res.agent) ?? null
      } catch {
        return null
      }
    },
    [refetch]
  )

  const deleteAgentApi = useCallback(
    async (agentId: string) => {
      await deleteAgent({ agent_id: agentId })
      await refetch()
    },
    [refetch]
  )

  useEffect(() => {
    refetch()
  }, [refetch])

  return {
    agents,
    loading,
    error,
    refetch,
    createAgent: createAgentApi,
    updateAgent: updateAgentApi,
    deleteAgent: deleteAgentApi,
  }
}
