import { useState, useCallback, useEffect } from 'react'
import { listProjects, createProject, deleteProject as deleteProjectApi, addIgnoredPath as addIgnoredPathApi, removeIgnoredPath as removeIgnoredPathApi } from '../api/client'
import { projectFromDto } from '../api/mappers'
import type { Project } from '../api/types'

export interface UseProjectsResult {
  projects: Project[]
  loading: boolean
  error: string | null
  refetch: () => Promise<void>
  createProject: (params: { name?: string; root_dir: string }) => Promise<Project | null>
  deleteProject: (projectId: string) => Promise<void>
  addIgnoredPath: (projectId: string, path: string) => Promise<void>
  removeIgnoredPath: (projectId: string, path: string) => Promise<void>
}

export function useProjects(): UseProjectsResult {
  const [projects, setProjects] = useState<Project[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refetch = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await listProjects()
      setProjects((res.projects ?? []).map(projectFromDto).filter(Boolean) as Project[])
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load projects')
      setProjects([])
    } finally {
      setLoading(false)
    }
  }, [])

  const createProjectApi = useCallback(
    async (params: { name?: string; root_dir: string }): Promise<Project | null> => {
      try {
        const res = await createProject({
          name: params.name ?? '',
          root_dir: params.root_dir,
        })
        await refetch()
        return projectFromDto(res.project) ?? null
      } catch {
        return null
      }
    },
    [refetch]
  )

  const deleteProjectCallback = useCallback(
    async (projectId: string) => {
      await deleteProjectApi({ project_id: projectId })
      await refetch()
    },
    [refetch]
  )

  const addIgnoredPathCallback = useCallback(
    async (projectId: string, path: string) => {
      await addIgnoredPathApi({ project_id: projectId, path })
      await refetch()
    },
    [refetch]
  )

  const removeIgnoredPathCallback = useCallback(
    async (projectId: string, path: string) => {
      await removeIgnoredPathApi({ project_id: projectId, path })
      await refetch()
    },
    [refetch]
  )

  useEffect(() => {
    refetch()
  }, [refetch])

  return {
    projects,
    loading,
    error,
    refetch,
    createProject: createProjectApi,
    deleteProject: deleteProjectCallback,
    addIgnoredPath: addIgnoredPathCallback,
    removeIgnoredPath: removeIgnoredPathCallback,
  }
}
