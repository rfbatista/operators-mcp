/** Tree node from list_tree tool */
export interface TreeNode {
  path: string
  name: string
  is_dir: boolean
  children?: TreeNode[]
}

/** Project from list_projects / get_project - defines the directory root for tree, paths, and zones */
export interface Project {
  id: string
  name: string
  root_dir: string
  ignored_paths: string[]
}

/** Zone from list_zones / get_zone */
export interface Zone {
  id: string
  project_id: string
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
