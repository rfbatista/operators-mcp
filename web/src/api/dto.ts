/**
 * Request/response DTOs matching the HTTP API (snake_case).
 * Used for serialization and mapping to/from UI models.
 */

/** Agent in zone DTO (and list/get/create/update agent) */
export interface AgentDto {
  id: string
  name: string
  description?: string
  prompt?: string
}

/** Zone DTO (API response shape) */
export interface ZoneDto {
  id: string
  project_id: string
  name: string
  pattern: string
  purpose: string
  constraints: string[]
  assigned_agents: AgentDto[]
  explicit_paths: string[]
}

/** Tree node DTO (API response shape) */
export interface TreeNodeDto {
  path: string
  name: string
  is_dir: boolean
  children: TreeNodeDto[]
}

/** Response: list_tree */
export interface ListTreeResponseDto {
  tree: TreeNodeDto
}

/** Project DTO (API response shape) */
export interface ProjectDto {
  id: string
  name: string
  root_dir: string
  ignored_paths?: string[]
}

/** Response: list_projects */
export interface ListProjectsResponseDto {
  projects: ProjectDto[]
}

/** Request: list_tree (query or body) */
export interface ListTreeRequestDto {
  root?: string
  project_id?: string
  depth?: number
}

/** Response: list_zones */
export interface ListZonesResponseDto {
  zones: ZoneDto[]
}

/** Request: list_zones */
export interface ListZonesRequestDto {
  project_id: string
}

/** Response: list_matching_paths */
export interface ListMatchingPathsResponseDto {
  paths: string[]
}

/** Request: list_matching_paths */
export interface ListMatchingPathsRequestDto {
  pattern: string
  root?: string
  project_id?: string
}

/** Request: create_project */
export interface CreateProjectRequestDto {
  name?: string
  root_dir: string
}

/** Response: create_project */
export interface CreateProjectResponseDto {
  project: ProjectDto
}

/** Request: delete_project */
export interface DeleteProjectRequestDto {
  project_id: string
}

/** Response: delete_project (204 No Content, or error) */

/** Request: add_ignored_path */
export interface AddIgnoredPathRequestDto {
  project_id: string
  path: string
}

/** Response: add_ignored_path */
export interface AddIgnoredPathResponseDto {
  project: ProjectDto
}

/** Request: remove_ignored_path */
export interface RemoveIgnoredPathRequestDto {
  project_id: string
  path: string
}

/** Response: remove_ignored_path */
export interface RemoveIgnoredPathResponseDto {
  project: ProjectDto
}

/** Response: get_zone */
export interface GetZoneResponseDto {
  zone: ZoneDto | null
}

/** Request: get_zone */
export interface GetZoneRequestDto {
  zone_id: string
}

/** Request: create_zone */
export interface CreateZoneRequestDto {
  project_id: string
  name: string
  pattern?: string
  purpose?: string
  constraints?: string[]
  assigned_agents?: AgentDto[]
}

/** Response: create_zone */
export interface CreateZoneResponseDto {
  zone: ZoneDto
}

/** Request: update_zone */
export interface UpdateZoneRequestDto {
  zone_id: string
  name?: string
  pattern?: string
  purpose?: string
  constraints?: string[]
  assigned_agents?: AgentDto[]
}

/** Response: update_zone */
export interface UpdateZoneResponseDto {
  zone: ZoneDto
}

/** Request: assign_path_to_zone */
export interface AssignPathToZoneRequestDto {
  zone_id: string
  path: string
}

/** Response: assign_path_to_zone */
export interface AssignPathToZoneResponseDto {
  zone: ZoneDto
}

/** Response: list_agents */
export interface ListAgentsResponseDto {
  agents: AgentDto[]
}

/** Request: get_agent */
export interface GetAgentRequestDto {
  agent_id: string
}

/** Response: get_agent */
export interface GetAgentResponseDto {
  agent: AgentDto
}

/** Request: create_agent */
export interface CreateAgentRequestDto {
  name?: string
  description?: string
  prompt?: string
}

/** Response: create_agent */
export interface CreateAgentResponseDto {
  agent: AgentDto
}

/** Request: update_agent */
export interface UpdateAgentRequestDto {
  agent_id: string
  name?: string
  description?: string
  prompt?: string
}

/** Response: update_agent */
export interface UpdateAgentResponseDto {
  agent: AgentDto
}

/** Request: delete_agent */
export interface DeleteAgentRequestDto {
  agent_id: string
}

/** API error response */
export interface ApiErrorDto {
  error: string
}
