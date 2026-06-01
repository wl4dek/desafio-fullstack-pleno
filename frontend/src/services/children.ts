import { api } from "@/lib/api"
import type { Alert, PaginatedResponse, Summary, ChildById } from "@/types"

export function fetchChildren(params: Record<string, string>) {
  const qs = new URLSearchParams(params).toString()
  return api.get<PaginatedResponse>(`/api/v1/children?${qs}`)
}

export function fetchChild(id: string) {
  return api.get<ChildById>(`/api/v1/children/${id}`)
}

export function fetchSummary() {
  return api.get<Summary>("/api/v1/summary")
}

export function markReviewed(id: string) {
  return api.patch<{ message: string }>(`/api/v1/children/${id}/review`)
}

export function fetchChildAlerts(id: string) {
  return api.get<Alert[]>(`/api/v1/children/${id}/alerts`)
}

export function listNeighborhood() {
  return api.get<string[]>(`/api/v1/children/neighborhood`)
}