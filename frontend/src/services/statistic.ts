import { api } from "@/lib/api"
import type { StatisticsResponse } from "@/types"

export function fetchStatistic(params: Record<string, string>) {
  const qs = new URLSearchParams(params).toString()
  return api.get<StatisticsResponse>(`/api/v1/statistics?${qs}`)
}