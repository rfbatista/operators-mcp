import { useState } from 'react'
import type { Project } from '../api/types'

export interface ProjectSelectorProps {
  projects: Project[]
  selectedId: string | null
  onSelect: (projectId: string | null) => void
  onCreateProject: (params: { name?: string; root_dir: string }) => Promise<Project | null>
  onDeleteProject?: (projectId: string) => Promise<void>
  loading?: boolean
  disabled?: boolean
}

export function ProjectSelector({
  projects,
  selectedId,
  onSelect,
  onCreateProject,
  onDeleteProject,
  loading,
  disabled,
}: ProjectSelectorProps) {
  const [showNewForm, setShowNewForm] = useState(false)
  const [newName, setNewName] = useState('')
  const [newRootDir, setNewRootDir] = useState('')
  const [creating, setCreating] = useState(false)
  const [createError, setCreateError] = useState<string | null>(null)
  const [deletingId, setDeletingId] = useState<string | null>(null)
  const [deleteConfirmId, setDeleteConfirmId] = useState<string | null>(null)

  const handleCreate = async () => {
    const root = newRootDir.trim()
    if (!root) {
      setCreateError('Root directory is required')
      return
    }
    setCreating(true)
    setCreateError(null)
    try {
      const created = await onCreateProject({
        name: newName.trim() || undefined,
        root_dir: root,
      })
      if (created) {
        onSelect(created.id)
        setShowNewForm(false)
        setNewName('')
        setNewRootDir('')
      } else {
        setCreateError('Failed to create project')
      }
    } catch (e) {
      setCreateError(e instanceof Error ? e.message : 'Failed to create project')
    } finally {
      setCreating(false)
    }
  }

  const handleDelete = async (projectId: string) => {
    if (!onDeleteProject) return
    setDeletingId(projectId)
    setDeleteConfirmId(null)
    try {
      await onDeleteProject(projectId)
      if (selectedId === projectId) {
        onSelect(null)
      }
    } finally {
      setDeletingId(null)
    }
  }

  return (
    <div className="flex flex-wrap items-center gap-3">
      <label className="flex items-center gap-2">
        <span className="font-semibold text-base-content">Project</span>
        <select
          value={selectedId ?? ''}
          onChange={(e) => onSelect(e.target.value === '' ? null : e.target.value)}
          disabled={disabled || loading}
          className="select select-bordered min-w-[200px]"
          aria-label="Select project"
        >
          <option value="">— Select project —</option>
          {projects.map((p) => (
            <option key={p.id} value={p.id}>
              {p.name || p.root_dir || p.id}
            </option>
          ))}
        </select>
      </label>
      {loading && <span className="text-sm text-base-content/60">Loading…</span>}
      <button
        type="button"
        className="btn btn-outline btn-sm"
        onClick={() => setShowNewForm((v) => !v)}
        disabled={disabled}
      >
        {showNewForm ? 'Cancel' : '+ New project'}
      </button>
      {onDeleteProject && selectedId && (
        <>
          {deleteConfirmId === selectedId ? (
            <span className="flex items-center gap-1 text-sm">
              <button
                type="button"
                className="btn btn-error btn-sm"
                onClick={() => handleDelete(selectedId)}
                disabled={deletingId !== null}
              >
                {deletingId === selectedId ? 'Deleting…' : 'Confirm delete'}
              </button>
              <button
                type="button"
                className="btn btn-ghost btn-sm"
                onClick={() => setDeleteConfirmId(null)}
                disabled={deletingId !== null}
              >
                Cancel
              </button>
            </span>
          ) : (
            <button
              type="button"
              className="btn btn-ghost btn-sm text-error"
              onClick={() => setDeleteConfirmId(selectedId)}
              disabled={disabled}
              title="Delete this project"
            >
              Delete project
            </button>
          )}
        </>
      )}

      {showNewForm && (
        <div className="card card-border bg-base-100 mt-2 min-w-[280px]">
          <div className="card-body gap-2 p-4">
            <h3 className="card-title text-sm">Create project</h3>
            {createError && (
              <div role="alert" className="alert alert-error alert-sm">
                <span>{createError}</span>
              </div>
            )}
            <label className="form-control w-full">
              <div className="label">
                <span className="label-text text-base-content/70">Name (optional)</span>
              </div>
              <input
                type="text"
                value={newName}
                onChange={(e) => setNewName(e.target.value)}
                placeholder="My project"
                className="input input-bordered input-sm w-full"
              />
            </label>
            <label className="form-control w-full">
              <div className="label">
                <span className="label-text text-base-content/70">Root directory *</span>
              </div>
              <input
                type="text"
                value={newRootDir}
                onChange={(e) => setNewRootDir(e.target.value)}
                placeholder="/path/to/project"
                className="input input-bordered input-sm w-full font-mono"
              />
            </label>
            <div className="card-actions mt-1 gap-2">
              <button
                type="button"
                className="btn btn-primary btn-sm"
                onClick={handleCreate}
                disabled={creating || !newRootDir.trim()}
              >
                {creating ? 'Creating…' : 'Create'}
              </button>
              <button
                type="button"
                className="btn btn-ghost btn-sm"
                onClick={() => {
                  setShowNewForm(false)
                  setCreateError(null)
                  setNewName('')
                  setNewRootDir('')
                }}
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
