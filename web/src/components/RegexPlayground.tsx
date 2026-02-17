import { useMatchingPaths } from '../hooks/useMatchingPaths'

export interface RegexPlaygroundProps {
  projectId: string | null
}

export function RegexPlayground({ projectId }: RegexPlaygroundProps) {
  const {
    pattern,
    setPattern,
    paths,
    loading,
    error,
    invalidPattern,
  } = useMatchingPaths(projectId)

  return (
    <section className="card card-border bg-base-100 mt-6">
      <div className="card-body">
        <h2 className="card-title text-lg">Regex playground</h2>
        <p className="text-sm text-base-content/70">
          Type a regex to see which paths match (debounced).
        </p>
        <div className="form-control w-full max-w-md">
          <input
            type="text"
            value={pattern}
            onChange={(e) => setPattern(e.target.value)}
            placeholder="e.g. cmd/.* or internal/"
            className={`input input-bordered w-full font-mono ${invalidPattern ? 'input-error' : ''}`}
            aria-invalid={invalidPattern}
          />
        </div>
        {invalidPattern && (
          <p className="text-sm text-error">Invalid pattern</p>
        )}
        {error && !invalidPattern && (
          <div role="alert" className="alert alert-error alert-sm">
            <span>{error}</span>
          </div>
        )}
        {loading && <p className="text-sm text-base-content/60">Loading…</p>}
        <div className="mt-2">
          <strong className="text-base-content">Matching paths</strong> ({paths.length}):
          <ul className="max-h-48 overflow-auto font-mono text-sm">
            {paths.length === 0 && !loading && pattern.trim() && !invalidPattern && (
              <li className="text-base-content/70">No matches</li>
            )}
            {paths.map((p) => (
              <li key={p}>{p}</li>
            ))}
          </ul>
        </div>
      </div>
    </section>
  )
}
