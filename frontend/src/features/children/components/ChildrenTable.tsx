"use client"

import { useRouter } from "next/navigation"
import { useChildren } from "@/hooks/useChildren"
import { AlertCategoryType, AlertsCategories, type ChildFilters } from "@/types"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton"
import { Button } from "@/components/ui/button"
import { Eye, RefreshCw } from "lucide-react"

interface ChildrenTableProps {
  filters: ChildFilters
}

export function ChildrenTable({ filters }: ChildrenTableProps) {
  const { data, isLoading, isError, refresh } = useChildren(filters)
  const router = useRouter()

  if (isLoading) {
    return (
      <div className="space-y-2">
        {Array.from({ length: 5 }).map((_, i) => (
          <Skeleton key={i} className="h-12 w-full" />
        ))}
      </div>
    )
  }

  if (isError) {
    return (
      <div className="text-center py-8">
        <p className="text-red-500 mb-2">Erro ao carregar crianças</p>
        <Button variant="outline" size="sm" onClick={() => refresh()} className="gap-2">
          <RefreshCw className="h-4 w-4" />
          Tentar novamente
        </Button>
      </div>
    )
  }

  if (!data || data.data.length === 0) {
    return (
      <div className="text-center py-12 text-neutral-500">
        <p>Nenhuma criança encontrada</p>
      </div>
    )
  }

  return (
    <div className="overflow-x-auto">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Nome</TableHead>
            <TableHead className="hidden md:table-cell">Bairro</TableHead>
            <TableHead className="text-center">Alerta</TableHead>
            <TableHead className="text-center">Revisão</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {data.data.map((child) => (
            <TableRow key={child.id}>
              <TableCell className="font-medium hover:cursor-pointer hover:bg-gray-600 dark:hover:bg-gray-400" onClick={() => router.push(`/children/${child.id}`)}>
                <div className="flex items-center gap-2">
                  <Eye className="h-4 w-4" />
                  {child.name}
                </div>
              </TableCell>
              <TableCell className="hidden md:table-cell">{child.neighborhood}</TableCell>
              <TableCell className="text-center pt-px">
                {child.alert_categories.length > 0 ? child.alert_categories.map((category, i) => (
                  <Badge className="mb-1 mr-1" key={i} variant="destructive">{AlertsCategories[category as AlertCategoryType]}</Badge>
                )) : <Badge className="mpx" variant="success">OK</Badge>}
              </TableCell>
              <TableCell className="text-center">
                {child.reviewed ? (
                  <Badge variant="success">Revisado</Badge>
                ) : (
                  <Badge variant="warning">Pendente</Badge>
                )}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
