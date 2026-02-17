import { useMatchingPaths } from '../hooks/useMatchingPaths'

export function RegexPlayground() {
  const {
    pattern,
    setPattern,
    paths,
    loading,
    error,
    invalidPattern,
  } = useMatchingPaths()

  return (
    <section style={{ marginTop: 16 }}>
      <h2>Regex playground</h2>
      <p style={{ fontSize: 12, color: '#666' }}>
        Type a regex to see which paths match (debounced).
      </p>
      <input
        type="text"
        value={pattern}
        onChange={(e) => setPattern(e.target.value)}
        placeholder="e.g. cmd/.* or internal/"
        style={{
          width: '100%',
          maxWidth: 400,
          padding: 8,
          fontFamily: 'monospace',
          border: invalidPattern ? '1px solid #c00' : '1px solid #ccc',
          borderRadius: 4,
        }}
        aria-invalid={invalidPattern}
      />
      {invalidPattern && (
        <p style={{ color: '#c00', fontSize: 12 }}>Invalid pattern</p>
      )}
      {error && !invalidPattern && (
        <p style={{ color: '#c00', fontSize: 12 }}>{error}</p>
      )}
      {loading && <p style={{ fontSize: 12, color: '#666' }}>Loading…</p>}
      <div style={{ marginTop: 8 }}>
        <strong>Matching paths</strong> ({paths.length}):
        <ul style={{ fontFamily: 'monospace', fontSize: 12, maxHeight: 200, overflow: 'auto' }}>
          {paths.length === 0 && !loading && pattern.trim() && !invalidPattern && (
            <li>No matches</li>
          )}
          {paths.map((p) => (
            <li key={p}>{p}</li>
          ))}
        </ul>
      </div>
    </section>
  )
}
