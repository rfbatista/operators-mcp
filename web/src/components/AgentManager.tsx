import { useState } from 'react'
import type { Agent } from '../api/types'

export interface AgentManagerProps {
  agents: Agent[]
  loading?: boolean
  onCreateAgent: (params: { name?: string; description?: string; prompt?: string }) => Promise<Agent | null>
  onUpdateAgent: (agentId: string, params: { name?: string; description?: string; prompt?: string }) => Promise<Agent | null>
  onDeleteAgent: (agentId: string) => Promise<void>
}

export function AgentManager({
  agents,
  loading,
  onCreateAgent,
  onUpdateAgent,
  onDeleteAgent,
}: AgentManagerProps) {
  const [showNewForm, setShowNewForm] = useState(false)
  const [newName, setNewName] = useState('')
  const [newDescription, setNewDescription] = useState('')
  const [newPrompt, setNewPrompt] = useState('')
  const [creating, setCreating] = useState(false)
  const [createError, setCreateError] = useState<string | null>(null)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [editName, setEditName] = useState('')
  const [editDescription, setEditDescription] = useState('')
  const [editPrompt, setEditPrompt] = useState('')
  const [saving, setSaving] = useState(false)
  const [deleteConfirmId, setDeleteConfirmId] = useState<string | null>(null)
  const [deletingId, setDeletingId] = useState<string | null>(null)

  const handleCreate = async () => {
    setCreating(true)
    setCreateError(null)
    try {
      const created = await onCreateAgent({
        name: newName.trim() || undefined,
        description: newDescription.trim() || undefined,
        prompt: newPrompt.trim() || undefined,
      })
      if (created) {
        setShowNewForm(false)
        setNewName('')
        setNewDescription('')
        setNewPrompt('')
      } else {
        setCreateError('Failed to create agent')
      }
    } catch (e) {
      setCreateError(e instanceof Error ? e.message : 'Failed to create agent')
    } finally {
      setCreating(false)
    }
  }

  const startEdit = (agent: Agent) => {
    setEditingId(agent.id)
    setEditName(agent.name)
    setEditDescription(agent.description)
    setEditPrompt(agent.prompt)
  }

  const handleUpdate = async () => {
    if (!editingId) return
    setSaving(true)
    try {
      await onUpdateAgent(editingId, {
        name: editName.trim() || undefined,
        description: editDescription.trim() || undefined,
        prompt: editPrompt.trim() || undefined,
      })
      setEditingId(null)
      setEditName('')
      setEditDescription('')
      setEditPrompt('')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (agentId: string) => {
    setDeletingId(agentId)
    setDeleteConfirmId(null)
    try {
      await onDeleteAgent(agentId)
    } finally {
      setDeletingId(null)
    }
  }

  return (
    <section className="card card-border bg-base-100 mt-4">
      <div className="card-body">
        <div className="flex flex-wrap items-center gap-2">
          <h2 className="card-title text-lg">Agents</h2>
          <button
            type="button"
            className="btn btn-outline btn-sm"
            onClick={() => setShowNewForm((v) => !v)}
            disabled={loading}
          >
            {showNewForm ? 'Cancel' : '+ New agent'}
          </button>
          {loading && (
            <span className="text-sm text-base-content/60">Loading…</span>
          )}
        </div>
        <p className="text-sm text-base-content/70">
          Agents can be assigned to zones. Create agents here, then assign them when editing a zone.
        </p>

        {showNewForm && (
          <div className="rounded-lg bg-base-200/50 p-3 space-y-3">
            <label className="form-control w-full">
              <span className="label-text text-base-content/70">Name</span>
              <input
                type="text"
                value={newName}
                onChange={(e) => setNewName(e.target.value)}
                placeholder="e.g. Backend Agent"
                className="input input-bordered input-sm w-full"
              />
            </label>
            <label className="form-control w-full">
              <span className="label-text text-base-content/70">Description</span>
              <textarea
                value={newDescription}
                onChange={(e) => setNewDescription(e.target.value)}
                placeholder="Short description of what this agent does"
                className="textarea textarea-bordered textarea-sm w-full min-h-[60px]"
                rows={2}
              />
            </label>
            <label className="form-control w-full">
              <span className="label-text text-base-content/70">Prompt</span>
              <textarea
                value={newPrompt}
                onChange={(e) => setNewPrompt(e.target.value)}
                placeholder="Instructions or system prompt for the agent"
                className="textarea textarea-bordered textarea-sm w-full min-h-[80px] font-mono text-sm"
                rows={3}
              />
            </label>
            <div className="flex flex-wrap items-center gap-2">
              <button
                type="button"
                className="btn btn-primary btn-sm"
                onClick={handleCreate}
                disabled={creating}
              >
                {creating ? 'Creating…' : 'Create'}
              </button>
              {createError && (
                <span className="text-error text-sm">{createError}</span>
              )}
            </div>
          </div>
        )}

        <ul className="menu rounded-box bg-base-200/50 w-full max-w-md">
          {agents.length === 0 && !loading && (
            <li className="text-base-content/60">
              <span>No agents yet. Create one to assign to zones.</span>
            </li>
          )}
          {agents.map((agent) => (
            <li key={agent.id}>
              <div className="flex w-full flex-col gap-2">
                {editingId === agent.id ? (
                  <div className="space-y-3 rounded-lg bg-base-300/50 p-3">
                    <label className="form-control w-full">
                      <span className="label-text text-base-content/70">Name</span>
                      <input
                        type="text"
                        value={editName}
                        onChange={(e) => setEditName(e.target.value)}
                        className="input input-bordered input-sm w-full"
                        placeholder="Agent name"
                        onKeyDown={(e) => {
                          if (e.key === 'Escape') {
                            setEditingId(null)
                            setEditName('')
                            setEditDescription('')
                            setEditPrompt('')
                          }
                        }}
                        autoFocus
                      />
                    </label>
                    <label className="form-control w-full">
                      <span className="label-text text-base-content/70">Description</span>
                      <textarea
                        value={editDescription}
                        onChange={(e) => setEditDescription(e.target.value)}
                        className="textarea textarea-bordered textarea-sm w-full min-h-[60px]"
                        rows={2}
                      />
                    </label>
                    <label className="form-control w-full">
                      <span className="label-text text-base-content/70">Prompt</span>
                      <textarea
                        value={editPrompt}
                        onChange={(e) => setEditPrompt(e.target.value)}
                        className="textarea textarea-bordered textarea-sm w-full min-h-[80px] font-mono text-sm"
                        rows={3}
                      />
                    </label>
                    <div className="flex gap-2">
                      <button
                        type="button"
                        className="btn btn-primary btn-sm"
                        onClick={handleUpdate}
                        disabled={saving}
                      >
                        {saving ? 'Saving…' : 'Save'}
                      </button>
                      <button
                        type="button"
                        className="btn btn-ghost btn-sm"
                        onClick={() => {
                          setEditingId(null)
                          setEditName('')
                          setEditDescription('')
                          setEditPrompt('')
                        }}
                        disabled={saving}
                      >
                        Cancel
                      </button>
                    </div>
                  </div>
                ) : (
                  <div className="flex w-full items-center justify-between gap-2 flex-wrap">
                    <div className="min-w-0 flex-1">
                      <span className="font-medium truncate block" title={agent.id}>
                        {agent.name || '(no name)'}
                      </span>
                      {agent.description && (
                        <span className="text-sm text-base-content/60 line-clamp-1" title={agent.description}>
                          {agent.description}
                        </span>
                      )}
                    </div>
                    <span className="flex items-center gap-1 shrink-0">
                      {deleteConfirmId === agent.id ? (
                        <>
                          <button
                            type="button"
                            className="btn btn-error btn-sm"
                            onClick={() => handleDelete(agent.id)}
                            disabled={deletingId !== null}
                          >
                            {deletingId === agent.id ? 'Deleting…' : 'Confirm'}
                          </button>
                          <button
                            type="button"
                            className="btn btn-ghost btn-sm"
                            onClick={() => setDeleteConfirmId(null)}
                            disabled={deletingId !== null}
                          >
                            Cancel
                          </button>
                        </>
                      ) : (
                        <>
                          <button
                            type="button"
                            className="btn btn-ghost btn-sm"
                            onClick={() => startEdit(agent)}
                          >
                            Edit
                          </button>
                          <button
                            type="button"
                            className="btn btn-ghost btn-sm text-error"
                            onClick={() => setDeleteConfirmId(agent.id)}
                            title="Delete agent"
                          >
                            Delete
                          </button>
                        </>
                      )}
                    </span>
                  </div>
                )}
              </div>
            </li>
          ))}
        </ul>
      </div>
    </section>
  )
}
