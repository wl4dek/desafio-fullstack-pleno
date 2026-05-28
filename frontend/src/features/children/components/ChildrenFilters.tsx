"use client"

import { useRouter, useSearchParams } from "next/navigation"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Search } from "lucide-react"
import { useNeighborhoods } from "@/hooks/useChildren"
import { useEffect, useState } from "react"

export function ChildrenFilters() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { neighborhoods } = useNeighborhoods() || { neighborhoods: [] }
  const [search, setSearch] = useState(searchParams.get("childName") || "")

  useEffect(() => {
    if (!search) {
      updateFilter("childName", "")
      return
    }

    const timeout = setTimeout(() => {
      updateFilter("childName", search)
    }, 500)

    return () => clearTimeout(timeout)
  }, [search])

  const updateFilter = (key: string, value: string) => {
    const params = new URLSearchParams(searchParams.toString())
    if (value) {
      params.set(key, value)
    } else {
      params.delete(key)
    }
    params.set("page", "1")
    router.replace(`/children?${params.toString()}`)
  }

  return (
    <div className="flex flex-col sm:flex-row gap-3 mb-4">
      <div className="flex-1">
        <div className="relative">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-neutral-500" />
          <Input
            placeholder="Buscar por criança..."
            className="pl-8"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
      </div>

      <div className="w-full sm:w-40">
        <Select
          value={searchParams.get("neighborhood") || ""}
          onValueChange={(v) => updateFilter("neighborhood", v)}
        >
          <SelectTrigger>
            <SelectValue placeholder="Bairro" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Todos</SelectItem>
            {neighborhoods?.map((n) => (
              <SelectItem key={n} value={n}>
                {n}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div className="w-full sm:w-40">
        <Select
          value={searchParams.get("has_alert") || ""}
          onValueChange={(v) => updateFilter("has_alert", v)}
        >
          <SelectTrigger>
            <SelectValue placeholder="Alerta" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Todos</SelectItem>
            <SelectItem value="true">Com alerta</SelectItem>
            <SelectItem value="false">Sem alerta</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div className="w-full sm:w-40">
        <Select
          value={searchParams.get("reviewed") || ""}
          onValueChange={(v) => updateFilter("reviewed", v)}
        >
          <SelectTrigger>
            <SelectValue placeholder="Revisão" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Todos</SelectItem>
            <SelectItem value="true">Revisado</SelectItem>
            <SelectItem value="false">Pendente</SelectItem>
          </SelectContent>
        </Select>
      </div>
    </div>
  )
}
