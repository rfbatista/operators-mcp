/**
 * HTTP client for the blueprint API. All endpoints return DTOs.
 * Uses relative /api when served from same origin.
 */

import type {
  ListTreeRequestDto,
  ListTreeResponseDto,
  ListZonesResponseDto,
  ListZonesRequestDto,
  ListMatchingPathsRequestDto,
  ListMatchingPathsResponseDto,
  ListProjectsResponseDto,
  CreateProjectRequestDto,
  CreateProjectResponseDto,
  DeleteProjectRequestDto,
  AddIgnoredPathRequestDto,
  AddIgnoredPathResponseDto,
  RemoveIgnoredPathRequestDto,
  RemoveIgnoredPathResponseDto,
  GetZoneRequestDto,
  GetZoneResponseDto,
  CreateZoneRequestDto,
  CreateZoneResponseDto,
  UpdateZoneRequestDto,
  UpdateZoneResponseDto,
  AssignPathToZoneRequestDto,
  AssignPathToZoneResponseDto,
  ListAgentsResponseDto,
  GetAgentRequestDto,
  GetAgentResponseDto,
  CreateAgentRequestDto,
  CreateAgentResponseDto,
  UpdateAgentRequestDto,
  UpdateAgentResponseDto,
  DeleteAgentRequestDto,
  ApiErrorDto,
} from './dto'

const API_BASE = '/api'

/** Options for request(); body may be any JSON-serializable value. */
type RequestOptions = Omit<RequestInit, 'body'> & { body?: unknown }

async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { method = 'GET', body, ...rest } = options
  const url = path.startsWith('http') ? path : `${API_BASE}${path}`
  const bodySerialized =
    body != null && method !== 'GET' ? JSON.stringify(body) : undefined
  const init: RequestInit = {
    ...rest,
    method,
    headers: { 'Content-Type': 'application/json', ...rest.headers },
    body: bodySerialized,
  }
  const res = await fetch(url, init)
  const data =
    res.status === 204 ? {} : await res.json().catch(() => ({}))
  if (!res.ok) {
    const msg = (data as ApiErrorDto).error ?? res.statusText
    throw new ApiError(res.status, msg)
  }
  if (res.status === 204) {
    return undefined as T
  }
  return data as T
}

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

/** GET list_projects */
export async function listProjects(): Promise<ListProjectsResponseDto> {
  return request<ListProjectsResponseDto>('/list_projects')
}

/** POST create_project */
export async function createProject(
  body: CreateProjectRequestDto
): Promise<CreateProjectResponseDto> {
  return request<CreateProjectResponseDto>('/create_project', {
    method: 'POST',
    body,
  })
}

/** POST delete_project */
export async function deleteProject(
  body: DeleteProjectRequestDto
): Promise<void> {
  await request<void>('/delete_project', {
    method: 'POST',
    body,
  })
}

/** POST add_ignored_path */
export async function addIgnoredPath(
  body: AddIgnoredPathRequestDto
): Promise<AddIgnoredPathResponseDto> {
  return request<AddIgnoredPathResponseDto>('/add_ignored_path', {
    method: 'POST',
    body,
  })
}

/** POST remove_ignored_path */
export async function removeIgnoredPath(
  body: RemoveIgnoredPathRequestDto
): Promise<RemoveIgnoredPathResponseDto> {
  return request<RemoveIgnoredPathResponseDto>('/remove_ignored_path', {
    method: 'POST',
    body,
  })
}

/** GET list_tree (optional query: root, project_id) */
export async function listTree(
  req: ListTreeRequestDto = {}
): Promise<ListTreeResponseDto> {
  const params = new URLSearchParams()
  if (req.root != null && req.root !== '') params.set('root', req.root)
  if (req.project_id != null && req.project_id !== '') params.set('project_id', req.project_id)
  const q = params.toString()
  return request<ListTreeResponseDto>(`/list_tree${q ? `?${q}` : ''}`)
}

/** GET list_zones?project_id=... */
export async function listZones(
  req: ListZonesRequestDto
): Promise<ListZonesResponseDto> {
  const params = new URLSearchParams()
  params.set('project_id', req.project_id)
  return request<ListZonesResponseDto>(`/list_zones?${params.toString()}`)
}

/** GET list_matching_paths?pattern=... (optional: root, project_id) */
export async function listMatchingPaths(
  req: ListMatchingPathsRequestDto
): Promise<ListMatchingPathsResponseDto> {
  const params = new URLSearchParams()
  params.set('pattern', req.pattern)
  if (req.root != null && req.root !== '') params.set('root', req.root)
  if (req.project_id != null && req.project_id !== '') params.set('project_id', req.project_id)
  return request<ListMatchingPathsResponseDto>(
    `/list_matching_paths?${params.toString()}`
  )
}

/** GET get_zone?zone_id=... */
export async function getZone(
  req: GetZoneRequestDto
): Promise<GetZoneResponseDto> {
  const params = new URLSearchParams()
  params.set('zone_id', req.zone_id)
  return request<GetZoneResponseDto>(`/get_zone?${params.toString()}`)
}

/** POST create_zone */
export async function createZone(
  body: CreateZoneRequestDto
): Promise<CreateZoneResponseDto> {
  return request<CreateZoneResponseDto>('/create_zone', {
    method: 'POST',
    body,
  })
}

/** POST update_zone */
export async function updateZone(
  body: UpdateZoneRequestDto
): Promise<UpdateZoneResponseDto> {
  return request<UpdateZoneResponseDto>('/update_zone', {
    method: 'POST',
    body,
  })
}

/** POST assign_path_to_zone */
export async function assignPathToZone(
  body: AssignPathToZoneRequestDto
): Promise<AssignPathToZoneResponseDto> {
  return request<AssignPathToZoneResponseDto>('/assign_path_to_zone', {
    method: 'POST',
    body,
  })
}

/** GET list_agents */
export async function listAgents(): Promise<ListAgentsResponseDto> {
  return request<ListAgentsResponseDto>('/list_agents')
}

/** GET get_agent?agent_id=... */
export async function getAgent(
  req: GetAgentRequestDto
): Promise<GetAgentResponseDto> {
  const params = new URLSearchParams()
  params.set('agent_id', req.agent_id)
  return request<GetAgentResponseDto>(`/get_agent?${params.toString()}`)
}

/** POST create_agent */
export async function createAgent(
  body: CreateAgentRequestDto
): Promise<CreateAgentResponseDto> {
  return request<CreateAgentResponseDto>('/create_agent', {
    method: 'POST',
    body,
  })
}

/** POST update_agent */
export async function updateAgent(
  body: UpdateAgentRequestDto
): Promise<UpdateAgentResponseDto> {
  return request<UpdateAgentResponseDto>('/update_agent', {
    method: 'POST',
    body,
  })
}

/** POST delete_agent */
export async function deleteAgent(
  body: DeleteAgentRequestDto
): Promise<void> {
  await request<void>('/delete_agent', {
    method: 'POST',
    body,
  })
}
