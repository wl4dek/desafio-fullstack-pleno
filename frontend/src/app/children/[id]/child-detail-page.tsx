"use client"

import { use } from "react"
import { ChildDetail } from "@/features/children/components/ChildDetail"

export function ChildDetailPage({ id }: { id: string }) {
  return <ChildDetail id={id} />
}
