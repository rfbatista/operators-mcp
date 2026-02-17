/**
 * HTTP client for the blueprint API. All endpoints return DTOs.
 * Uses relative /api when served from same origin.
 */

import type {
  ListTreeRequestDto,
  ListTreeResponseDto,
  ListZonesResponseDto,
  ListMatchingPathsRequestDto,
  ListMatchingPathsResponseDto,
  GetZoneRequestDto,
  GetZoneResponseDto,
  CreateZoneRequestDto,
  CreateZoneResponseDto,
  UpdateZoneRequestDto,
  UpdateZoneResponseDto,
  AssignPathToZoneRequestDto,
  AssignPathToZoneResponseDto,
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
  const data = await res.json().catch(() => ({}))
  if (!res.ok) {
    const msg = (data as ApiErrorDto).error ?? res.statusText
    throw new ApiError(res.status, msg)
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

/** GET list_tree (optional query: root) */
export async function listTree(
  req: ListTreeRequestDto = {}
): Promise<ListTreeResponseDto> {
  const params = new URLSearchParams()
  if (req.root != null && req.root !== '') params.set('root', req.root)
  const q = params.toString()
  return request<ListTreeResponseDto>(`/list_tree${q ? `?${q}` : ''}`)
}

/** GET list_zones */
export async function listZones(): Promise<ListZonesResponseDto> {
  return request<ListZonesResponseDto>('/list_zones')
}

/** GET list_matching_paths?pattern=... (optional: root) */
export async function listMatchingPaths(
  req: ListMatchingPathsRequestDto
): Promise<ListMatchingPathsResponseDto> {
  const params = new URLSearchParams()
  params.set('pattern', req.pattern)
  if (req.root != null && req.root !== '') params.set('root', req.root)
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
