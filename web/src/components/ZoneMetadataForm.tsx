import { useState, useEffect } from 'react'
import type { Zone } from '../api/types'

export interface ZoneMetadataFormProps {
  zoneId: string | null
  initialZone: Zone | null
  onLoadZone: (id: string) => Promise<Zone | null>
  onCreateZone: (params: {
    name: string
    pattern?: string
    purpose?: string
    constraints?: string[]
    assigned_agent?: string
  }) => Promise<Zone | null>
  onUpdateZone: (params: {
    zone_id: string
    name?: string
    pattern?: string
    purpose?: string
    constraints?: string[]
    assigned_agent?: string
  }) => Promise<Zone | null>
  onSaved?: () => void
}

export function ZoneMetadataForm({
  zoneId,
  initialZone,
  onLoadZone,
  onCreateZone,
  onUpdateZone,
  onSaved,
}: ZoneMetadataFormProps) {
  const [name, setName] = useState('')
  const [pattern, setPattern] = useState('')
  const [purpose, setPurpose] = useState('')
  const [constraints, setConstraints] = useState('')
  const [assignedAgent, setAssignedAgent] = useState('')
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const isEdit = zoneId != null && zoneId !== ''

  useEffect(() => {
    if (initialZone) {
      setName(initialZone.name)
      setPattern(initialZone.pattern ?? '')
      setPurpose(initialZone.purpose ?? '')
      setConstraints((initialZone.constraints ?? []).join('\n'))
      setAssignedAgent(initialZone.assigned_agent ?? '')
    } else if (!isEdit) {
      setName('')
      setPattern('')
      setPurpose('')
      setConstraints('')
      setAssignedAgent('')
    }
  }, [initialZone, isEdit])

  useEffect(() => {
    if (zoneId && !initialZone) {
      onLoadZone(zoneId).then((z) => {
        if (z) {
          setName(z.name)
          setPattern(z.pattern ?? '')
          setPurpose(z.purpose ?? '')
          setConstraints((z.constraints ?? []).join('\n'))
          setAssignedAgent(z.assigned_agent ?? '')
        }
      })
    }
  }, [zoneId, initialZone, onLoadZone])

  const handleSave = async () => {
    setSaving(true)
    setError(null)
    try {
      const constraintList = constraints
        .split('\n')
        .map((s) => s.trim())
        .filter(Boolean)
      if (isEdit && zoneId) {
        await onUpdateZone({
          zone_id: zoneId,
          name,
          pattern,
          purpose,
          constraints: constraintList,
          assigned_agent: assignedAgent,
        })
      } else {
        await onCreateZone({
          name: name || 'Unnamed zone',
          pattern,
          purpose,
          constraints: constraintList,
          assigned_agent: assignedAgent,
        })
      }
      onSaved?.()
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to save')
    } finally {
      setSaving(false)
    }
  }

  return (
    <section className="card card-border bg-base-100 mt-6">
      <div className="card-body">
        <h2 className="card-title text-lg">{isEdit ? 'Edit zone' : 'New zone'}</h2>
        {error && (
          <div role="alert" className="alert alert-error">
            <span>{error}</span>
          </div>
        )}
        <div className="flex max-w-md flex-col gap-3">
          <label className="form-control w-full">
            <div className="label">
              <span className="label-text">Name</span>
            </div>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              className="input input-bordered w-full"
            />
          </label>
          <label className="form-control w-full">
            <div className="label">
              <span className="label-text">Pattern (regex)</span>
            </div>
            <input
              type="text"
              value={pattern}
              onChange={(e) => setPattern(e.target.value)}
              className="input input-bordered w-full font-mono"
            />
          </label>
          <label className="form-control w-full">
            <div className="label">
              <span className="label-text">Purpose</span>
            </div>
            <textarea
              value={purpose}
              onChange={(e) => setPurpose(e.target.value)}
              rows={2}
              className="textarea textarea-bordered w-full"
            />
          </label>
          <label className="form-control w-full">
            <div className="label">
              <span className="label-text">Constraints (one per line)</span>
            </div>
            <textarea
              value={constraints}
              onChange={(e) => setConstraints(e.target.value)}
              rows={3}
              className="textarea textarea-bordered w-full font-mono"
            />
          </label>
          <label className="form-control w-full">
            <div className="label">
              <span className="label-text">Assigned agent</span>
            </div>
            <input
              type="text"
              value={assignedAgent}
              onChange={(e) => setAssignedAgent(e.target.value)}
              className="input input-bordered w-full"
            />
          </label>
          <button
            type="button"
            className="btn btn-primary"
            onClick={handleSave}
            disabled={saving || !name.trim()}
          >
            {saving ? 'Saving…' : isEdit ? 'Update zone' : 'Create zone'}
          </button>
        </div>
      </div>
    </section>
  )
}
