import { useState } from 'react'
import type { Zone } from '../api/types'

export interface AssignToZoneProps {
  path: string
  zones: Zone[]
  onSelect: (zoneId: string) => void
  onCreateNew: () => void
  onCancel: () => void
}

export function AssignToZone({
  path,
  zones,
  onSelect,
  onCreateNew,
  onCancel,
}: AssignToZoneProps) {
  const [selectedId, setSelectedId] = useState<string>('')

  return (
    <div
      style={{
        position: 'fixed',
        inset: 0,
        background: 'rgba(0,0,0,0.3)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        zIndex: 1000,
      }}
      role="dialog"
      aria-label="Assign to zone"
    >
      <div
        style={{
          background: 'white',
          padding: 16,
          borderRadius: 8,
          minWidth: 280,
          boxShadow: '0 4px 12px rgba(0,0,0,0.15)',
        }}
      >
        <h3>Assign to zone</h3>
        <p style={{ fontSize: 12, color: '#666' }}>Path: {path}</p>
        <select
          value={selectedId}
          onChange={(e) => setSelectedId(e.target.value)}
          style={{ width: '100%', padding: 8, marginBottom: 8 }}
        >
          <option value="">Select a zone…</option>
          {zones.map((z) => (
            <option key={z.id} value={z.id}>
              {z.name}
            </option>
          ))}
        </select>
        <div style={{ display: 'flex', gap: 8, marginTop: 12 }}>
          <button
            type="button"
            onClick={() => selectedId && onSelect(selectedId)}
            disabled={!selectedId}
          >
            Assign
          </button>
          <button type="button" onClick={onCreateNew}>
            Create new zone
          </button>
          <button type="button" onClick={onCancel}>
            Cancel
          </button>
        </div>
      </div>
    </div>
  )
}
