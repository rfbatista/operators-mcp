/**
 * Maps API DTOs to UI models and UI params to request DTOs.
 */

import type { ZoneDto, AgentDto, TreeNodeDto } from './dto'
import type { Zone, TreeNode } from './types'

/** Map ZoneDto to UI Zone (assigned_agent = first assigned agent name) */
export function zoneFromDto(dto: ZoneDto | null | undefined): Zone | null {
  if (dto == null) return null
  const assigned_agent =
    dto.assigned_agents?.length > 0 ? dto.assigned_agents[0].name : ''
  return {
    id: dto.id,
    name: dto.name,
    pattern: dto.pattern ?? '',
    purpose: dto.purpose ?? '',
    constraints: dto.constraints ?? [],
    assigned_agent,
    explicit_paths: dto.explicit_paths ?? [],
  }
}

/** Map UI Zone or create params to assigned_agents DTO array */
export function toAssignedAgentsDto(assigned_agent?: string): AgentDto[] {
  if (assigned_agent == null || assigned_agent.trim() === '') return []
  return [{ id: '', name: assigned_agent.trim() }]
}

/** Map TreeNodeDto to UI TreeNode */
export function treeNodeFromDto(dto: TreeNodeDto | null | undefined): TreeNode | null {
  if (dto == null) return null
  return {
    path: dto.path,
    name: dto.name,
    is_dir: dto.is_dir,
    children: (dto.children ?? []).map(treeNodeFromDto).filter(Boolean) as TreeNode[],
  }
}
