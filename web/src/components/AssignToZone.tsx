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
    <dialog className="modal modal-open" open role="dialog" aria-label="Assign to zone">
      <div
        className="modal-backdrop bg-black/50"
        aria-hidden="true"
        onClick={onCancel}
        onKeyDown={(e) => e.key === 'Escape' && onCancel()}
      />
      <div className="modal-box">
        <h3 className="text-lg font-bold">Assign to zone</h3>
        <p className="text-sm text-base-content/70">Path: {path}</p>
        <div className="form-control mt-2 w-full">
          <select
            value={selectedId}
            onChange={(e) => setSelectedId(e.target.value)}
            className="select select-bordered w-full"
          >
            <option value="">Select a zone…</option>
            {zones.map((z) => (
              <option key={z.id} value={z.id}>
                {z.name}
              </option>
            ))}
          </select>
        </div>
        <div className="modal-action">
          <button
            type="button"
            className="btn btn-primary"
            onClick={() => selectedId && onSelect(selectedId)}
            disabled={!selectedId}
          >
            Assign
          </button>
          <button type="button" className="btn btn-outline" onClick={onCreateNew}>
            Create new zone
          </button>
          <button type="button" className="btn btn-ghost" onClick={onCancel}>
            Cancel
          </button>
        </div>
      </div>
    </dialog>
  )
}
