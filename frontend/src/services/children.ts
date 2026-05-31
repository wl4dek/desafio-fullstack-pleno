import { api } from "@/lib/api"
import type { Alert, PaginatedResponse, Summary, ChildById } from "@/types"

export function fetchChildren(params: Record<string, string>) {
  const qs = new URLSearchParams(params).toString()
  return api.get<PaginatedResponse>(`/children?${qs}`)
}

export function fetchChild(id: string) {
  return api.get<ChildById>(`/children/${id}`)
}

export function fetchSummary() {
  return api.get<Summary>("/summary")
}

export function markReviewed(id: string) {
  return api.patch<{ message: string }>(`/children/${id}/review`)
}

export function fetchChildAlerts(id: string) {
  return api.get<Alert[]>(`/children/${id}/alerts`)
}

export function listNeighborhood() {
  return api.get<string[]>(`/children/neighborhood`)
}