/** Tree node from list_tree tool */
export interface TreeNode {
  path: string
  name: string
  is_dir: boolean
  children?: TreeNode[]
}

/** Zone from list_zones / get_zone */
export interface Zone {
  id: string
  name: string
  pattern: string
  purpose: string
  constraints: string[]
  assigned_agent: string
  explicit_paths: string[]
}

/** Tool call result (content text is JSON) */
export interface CallToolResult {
  content?: Array<{ text?: string }>
  isError?: boolean
}

declare global {
  interface Window {
    __callTool__?: (name: string, args: Record<string, unknown>) => Promise<CallToolResult>
  }
}
