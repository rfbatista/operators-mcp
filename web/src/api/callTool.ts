import type { CallToolResult } from './types'

/**
 * Invoke an MCP tool. When the designer runs inside an MCP host (e.g. Cursor),
 * the host should set window.__callTool__ so the UI can call the backend.
 * When not set, returns mock/empty data so the UI still renders.
 */
export async function callTool(
  name: string,
  args: Record<string, unknown> = {}
): Promise<CallToolResult> {
  if (typeof window !== 'undefined' && window.__callTool__) {
    return window.__callTool__(name, args)
  }
  // Fallback: mock data for dev so the UI is usable without a host bridge
  return mockCallTool(name, args)
}

function mockCallTool(
  name: string,
  _args: Record<string, unknown>
): Promise<CallToolResult> {
  switch (name) {
    case 'list_tree':
      return Promise.resolve({
        content: [
          {
            text: JSON.stringify({
              tree: {
                path: '',
                name: '.',
                is_dir: true,
                children: [
                  { path: 'cmd', name: 'cmd', is_dir: true, children: [] },
                  { path: 'internal', name: 'internal', is_dir: true, children: [] },
                  { path: 'web', name: 'web', is_dir: true, children: [] },
                ],
              },
            }),
          },
        ],
      })
    case 'list_zones':
      return Promise.resolve({
        content: [{ text: JSON.stringify({ zones: [] }) }],
      })
    case 'list_matching_paths':
      return Promise.resolve({
        content: [{ text: JSON.stringify({ paths: [] }) }],
      })
    case 'get_zone':
      return Promise.resolve({ isError: true })
    case 'create_zone':
    case 'update_zone':
    case 'assign_path_to_zone':
      return Promise.resolve({ isError: true })
    default:
      return Promise.resolve({ isError: true })
  }
}
