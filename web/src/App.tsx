import { useState } from 'react'
import type { Zone } from './api/types'
import { useListTree } from './hooks/useListTree'
import { useZoneHighlights } from './hooks/useZoneHighlights'
import { useZones } from './hooks/useZones'
import { useProjects } from './hooks/useProjects'
import { SourceTree } from './components/SourceTree'
import { RegexPlayground } from './components/RegexPlayground'
import { ZoneMetadataForm, type ZoneMetadataFormProps } from './components/ZoneMetadataForm'
import { AssignToZone } from './components/AssignToZone'
import { ProjectSelector } from './components/ProjectSelector'
import { ThemeSwitcher } from './components/ThemeSwitcher'

export default function App() {
  const [selectedProjectId, setSelectedProjectId] = useState<string | null>(null)

  const {
    projects,
    loading: projectsLoading,
    error: projectsError,
    createProject,
    addIgnoredPath,
    removeIgnoredPath,
  } = useProjects()

  const { tree, loading: treeLoading, error: treeError, refetch: refetchTree } = useListTree(
    selectedProjectId
  )
  const {
    highlightPaths,
    pathToZones,
    zones: highlightZones,
    loading: zonesLoading,
    error: zonesError,
    refetch: refetchZones,
  } = useZoneHighlights(selectedProjectId)
  const {
    zones,
    refetch: refetchZonesList,
    getZone,
    createZone,
    updateZone,
    assignPathToZone,
  } = useZones(selectedProjectId)

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

  const handleCreateZone = async (params: {
    name: string
    pattern?: string
    purpose?: string
    constraints?: string[]
    assigned_agent?: string
  }) => {
    if (!selectedProjectId) return null
    return createZone({ ...params, project_id: selectedProjectId })
  }

  const loading = treeLoading || zonesLoading
  const error = treeError ?? zonesError ?? projectsError

  return (
    <div className="min-h-screen bg-base-200 p-4 md:p-6">
      <div className="mx-auto max-w-4xl">
        <header className="flex flex-wrap items-center justify-between gap-4">
          <h1 className="text-2xl font-bold text-base-content md:text-3xl">
            Architecture Designer
          </h1>
          <ThemeSwitcher />
        </header>

        <div className="mt-4">
          <ProjectSelector
            projects={projects}
            selectedId={selectedProjectId}
            onSelect={setSelectedProjectId}
            onCreateProject={createProject}
            loading={projectsLoading}
          />
        </div>

        {error && (
          <div role="alert" className="alert alert-error mt-4">
            <span>{error}</span>
          </div>
        )}

        {!selectedProjectId && !projectsLoading && (
          <p className="mt-4 text-base-content/70">
            Select a project above to view the source tree and zones, or create a new one.
          </p>
        )}

        {selectedProjectId && loading && !tree && (
          <p className="mt-4 text-base-content/70">Loading project tree…</p>
        )}
        {selectedProjectId && !tree && !loading && !error && (
          <p className="mt-4 text-base-content/70">
            No project tree available. Check that the root directory exists and is readable.
          </p>
        )}

        {tree && (
          <section className="card card-border bg-base-100 mt-6">
            <div className="card-body">
              <h2 className="card-title text-lg">Source tree</h2>
              <SourceTree
                tree={tree}
                highlightPaths={highlightPaths}
                pathToZones={pathToZones}
                ignoredPaths={selectedProjectId ? (projects.find((p) => p.id === selectedProjectId)?.ignored_paths ?? []) : []}
                onAssignToZone={handleAssignToZone}
                onHide={selectedProjectId ? (path) => addIgnoredPath(selectedProjectId, path) : undefined}
                onShow={selectedProjectId ? (path) => removeIgnoredPath(selectedProjectId, path) : undefined}
              />
            </div>
          </section>
        )}

        {selectedProjectId && <RegexPlayground projectId={selectedProjectId} />}

        {showNewZone && selectedProjectId && (
          <ZoneMetadataForm
            zoneId={null}
            initialZone={null}
            onLoadZone={getZone}
            onCreateZone={handleCreateZone as ZoneMetadataFormProps['onCreateZone']}
            onUpdateZone={updateZone}
            onSaved={handleZoneSaved}
          />
        )}

        {editingZoneId && (
          <ZoneMetadataForm
            zoneId={editingZoneId}
            initialZone={editingZone}
            onLoadZone={getZone}
            onCreateZone={handleCreateZone as ZoneMetadataFormProps['onCreateZone']}
            onUpdateZone={updateZone}
            onSaved={handleZoneSaved}
          />
        )}

        <section className="card card-border bg-base-100 mt-6">
          <div className="card-body">
            <div className="flex flex-wrap items-center gap-2">
              <h2 className="card-title text-lg">Zones</h2>
              <button
                type="button"
                className="btn btn-primary btn-sm"
                onClick={() => setShowNewZone(true)}
                disabled={!selectedProjectId}
              >
                New zone
              </button>
              {!selectedProjectId && (
                <span className="text-sm text-base-content/60">Select a project first</span>
              )}
            </div>
            <ul className="menu rounded-box bg-base-200/50 w-full max-w-md">
              {zones.map((z) => (
                <li key={z.id}>
                  <div className="flex w-full items-center justify-between gap-2">
                    <span className="font-medium">{z.name}</span>
                    <button
                      type="button"
                      className="btn btn-ghost btn-sm"
                      onClick={() => {
                        setEditingZoneId(z.id)
                        setEditingZone(z)
                      }}
                    >
                      Edit
                    </button>
                  </div>
                </li>
              ))}
            </ul>
          </div>
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
    </div>
  )
}
