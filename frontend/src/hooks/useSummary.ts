import useSWR from "swr"
import { fetchSummary } from "@/services/children"
import type { Summary } from "@/types"

export function useSummary() {
  const { data, error, isLoading, mutate } = useSWR<Summary>(
    "/summary",
    fetchSummary,
  )

  return {
    summary: data,
    isLoading,
    isError: !!error,
    error,
    refresh: mutate,
  }
}
