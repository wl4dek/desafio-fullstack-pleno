import useSWR from "swr"
import { fetchChildren } from "@/services/children"
import type { PaginatedResponse, ChildFilters, Alert } from "@/types"

export function useChildren(filters: ChildFilters) {
  const params = Object.fromEntries(
    Object.entries(filters)
      .map(([k, v]) => [k, v === 'all' ? '' : v])
      .filter(([, v]) => v !== undefined && v !== ""),
  )

  const qs = new URLSearchParams(params).toString()
  const key = `/children?${qs}`

  const { data, error, isLoading, mutate } = useSWR<PaginatedResponse>(
    key,
    () => fetchChildren(params),
  )

  return {
    data,
    isLoading,
    isError: !!error,
    error,
    refresh: mutate,
  }
}

export function useChild(id: string) {
  const { data, error, isLoading } = useSWR(
    id ? `/children/${id}` : null,
    () => import("@/services/children").then((m) => m.fetchChild(id)),
  )

  return {
    child: data,
    isLoading,
    isError: !!error,
    error,
  }
}

export function useNeighborhoods() {
  const { data, error, isLoading } = useSWR(
    `/children/neighborhood`,
    () => import("@/services/children").then((m) => m.listNeighborhood()),
  )

  return {
    neighborhoods: data,
    isLoading,
    isError: !!error,
    error,
  }
}