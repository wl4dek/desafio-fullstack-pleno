"use client"

import { use, useCallback } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { ChildrenFilters } from "@/features/children/components/ChildrenFilters"
import { ChildrenTable } from "@/features/children/components/ChildrenTable"
import { Button } from "@/components/ui/button"
import { ChevronLeft, ChevronRight } from "lucide-react"

export function ChildrenList() {
  const searchParams = useSearchParams()
  const router = useRouter()

  const filters = {
    childName: searchParams.get("childName") || "",
    neighborhood: searchParams.get("neighborhood") || "",
    alert: searchParams.get("alert") || "",
    has_alert: searchParams.get("has_alert") || "",
    reviewed: searchParams.get("reviewed") || "",
    page: searchParams.get("page") || "1",
    per_page: searchParams.get("per_page") || "10",
  }

  const page = parseInt(filters.page, 10)

  const goToPage = useCallback(
    (p: number) => {
      const params = new URLSearchParams(searchParams.toString())
      params.set("page", String(p))
      router.push(`/children?${params.toString()}`)
    },
    [searchParams, router],
  )

  return (
    <div>
      <ChildrenFilters />
      <ChildrenTable filters={filters} />
      <div className="flex items-center justify-between mt-4">
        <p className="text-sm text-neutral-500">
          Página {page}
        </p>
        <div className="flex gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => goToPage(page - 1)}
            disabled={page <= 1}
          >
            <ChevronLeft className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => goToPage(page + 1)}
          >
            <ChevronRight className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  )
}
