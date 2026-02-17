/**
 * Maps API DTOs to UI models and UI params to request DTOs.
 */

import type { ZoneDto, AgentDto, TreeNodeDto, ProjectDto } from './dto'
import type { Zone, TreeNode, Project, Agent } from './types'

/** Map AgentDto to UI Agent */
export function agentFromDto(dto: AgentDto | null | undefined): Agent | null {
  if (dto == null) return null
  return {
    id: dto.id,
    name: dto.name ?? '',
    description: dto.description ?? '',
    prompt: dto.prompt ?? '',
  }
}

/** Map ProjectDto to UI Project */
export function projectFromDto(dto: ProjectDto | null | undefined): Project | null {
  if (dto == null) return null
  return {
    id: dto.id,
    name: dto.name ?? '',
    root_dir: dto.root_dir ?? '',
    ignored_paths: dto.ignored_paths ?? [],
  }
}

/** Map ZoneDto to UI Zone (assigned_agent/assigned_agent_id from first assigned agent) */
export function zoneFromDto(dto: ZoneDto | null | undefined): Zone | null {
  if (dto == null) return null
  const first = dto.assigned_agents?.[0]
  return {
    id: dto.id,
    project_id: dto.project_id ?? '',
    name: dto.name,
    pattern: dto.pattern ?? '',
    purpose: dto.purpose ?? '',
    constraints: dto.constraints ?? [],
    assigned_agent: first?.name ?? '',
    assigned_agent_id: first?.id ?? '',
    explicit_paths: dto.explicit_paths ?? [],
  }
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
