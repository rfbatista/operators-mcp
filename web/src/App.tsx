import { useState } from 'react'
import type { Zone } from './api/types'
import { useListTree } from './hooks/useListTree'
import { useZoneHighlights } from './hooks/useZoneHighlights'
import { useZones } from './hooks/useZones'
import { SourceTree } from './components/SourceTree'
import { RegexPlayground } from './components/RegexPlayground'
import { ZoneMetadataForm } from './components/ZoneMetadataForm'
import { AssignToZone } from './components/AssignToZone'

export default function App() {
  const { tree, loading: treeLoading, error: treeError, refetch: refetchTree } = useListTree()
  const {
    highlightPaths,
    pathToZones,
    zones: highlightZones,
    loading: zonesLoading,
    error: zonesError,
    refetch: refetchZones,
  } = useZoneHighlights()
  const {
    zones,
    refetch: refetchZonesList,
    getZone,
    createZone,
    updateZone,
    assignPathToZone,
  } = useZones()

  const [assignPath, setAssignPath] = useState<string | null>(null)
  const [editingZoneId, setEditingZoneId] = useState<string | null>(null)
  const [editingZone, setEditingZone] = useState<Zone | null>(null)
  const [showNewZone, setShowNewZone] = useState(false)

  const handleAssignToZone = (path: string) => {
    setAssignPath(path)
  }

  const handleAssignSelect = async (zoneId: string) => {
    if (!assignPath) return
    await assignPathToZone(zoneId, assignPath)
    setAssignPath(null)
    refetchZones()
    refetchTree()
    refetchZonesList()
  }

  const handleAssignCreateNew = () => {
    setShowNewZone(true)
    setAssignPath(null)
  }

  const handleZoneSaved = () => {
    setEditingZoneId(null)
    setEditingZone(null)
    setShowNewZone(false)
    refetchZones()
    refetchZonesList()
    refetchTree()
  }

  const loading = treeLoading || zonesLoading
  const error = treeError ?? zonesError

  return (
    <div style={{ padding: 16, maxWidth: 900 }}>
      <h1>Architecture Designer</h1>
      {error && <p style={{ color: '#c00' }}>{error}</p>}
      {loading && !tree && <p>Loading project tree…</p>}
      {!tree && !loading && !error && (
        <p>No project tree available. The host may not have provided a backend bridge.</p>
      )}

      {tree && (
        <section>
          <h2>Source tree</h2>
          <SourceTree
            tree={tree}
            highlightPaths={highlightPaths}
            pathToZones={pathToZones}
            onAssignToZone={handleAssignToZone}
          />
        </section>
      )}

      <RegexPlayground />

      {showNewZone && (
        <ZoneMetadataForm
          zoneId={null}
          initialZone={null}
          onLoadZone={getZone}
          onCreateZone={createZone}
          onUpdateZone={updateZone}
          onSaved={handleZoneSaved}
        />
      )}

      {editingZoneId && (
        <ZoneMetadataForm
          zoneId={editingZoneId}
          initialZone={editingZone}
          onLoadZone={getZone}
          onCreateZone={createZone}
          onUpdateZone={updateZone}
          onSaved={handleZoneSaved}
        />
      )}

      <section style={{ marginTop: 16 }}>
        <h2>Zones</h2>
        <button type="button" onClick={() => setShowNewZone(true)}>
          New zone
        </button>
        <ul>
          {zones.map((z) => (
            <li key={z.id}>
              <strong>{z.name}</strong>
              <button
                type="button"
                onClick={() => {
                  setEditingZoneId(z.id)
                  setEditingZone(z)
                }}
                style={{ marginLeft: 8 }}
              >
                Edit
              </button>
            </li>
          ))}
        </ul>
      </section>

      {assignPath && (
        <AssignToZone
          path={assignPath}
          zones={highlightZones.length > 0 ? highlightZones : zones}
          onSelect={handleAssignSelect}
          onCreateNew={handleAssignCreateNew}
          onCancel={() => setAssignPath(null)}
        />
      )}
    </div>
  )
}
