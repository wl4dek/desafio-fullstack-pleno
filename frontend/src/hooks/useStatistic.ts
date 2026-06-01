import useSWR from "swr"
import { fetchStatistic } from "@/services/statistic"
import type { StatisticsResponse } from "@/types"

export function useStatistic() {
  const { data, error, isLoading, mutate } = useSWR<StatisticsResponse>(
    "/api/v1/statistics",
    fetchStatistic,
  )

  return {
    statistics: data?.statistics ?? [],
    isLoading,
    isError: !!error,
    error,
    refresh: mutate,
  }
}
