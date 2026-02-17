/**
 * Request/response DTOs matching the HTTP API (snake_case).
 * Used for serialization and mapping to/from UI models.
 */

/** Agent in zone DTO */
export interface AgentDto {
  id: string
  name: string
}

/** Zone DTO (API response shape) */
export interface ZoneDto {
  id: string
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

/** Request: list_tree (query or body) */
export interface ListTreeRequestDto {
  root?: string
  depth?: number
}

/** Response: list_zones */
export interface ListZonesResponseDto {
  zones: ZoneDto[]
}

/** Response: list_matching_paths */
export interface ListMatchingPathsResponseDto {
  paths: string[]
}

/** Request: list_matching_paths */
export interface ListMatchingPathsRequestDto {
  pattern: string
  root?: string
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

/** API error response */
export interface ApiErrorDto {
  error: string
}
