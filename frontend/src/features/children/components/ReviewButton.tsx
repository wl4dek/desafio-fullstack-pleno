"use client"

import { useState } from "react"
import { useSWRConfig } from "swr"
import { markReviewed } from "@/services/children"
import { Button } from "@/components/ui/button"
import { toast } from "@/hooks/use-toast"
import { CheckCircle } from "lucide-react"

interface ReviewButtonProps {
  childId: string
}

export function ReviewButton({ childId }: ReviewButtonProps) {
  const [loading, setLoading] = useState(false)
  const { mutate } = useSWRConfig()

  const handleReview = async () => {
    setLoading(true)
    try {
      await markReviewed(childId)
      await mutate(`/children/${childId}`)
      await mutate((key: string) => typeof key === "string" && key.startsWith("/children?"))
      await mutate("/summary")
      toast({
        title: "Revisão registrada",
        description: "Criança marcada como revisada com sucesso.",
        variant: "success",
      })
    } catch {
      toast({
        title: "Erro",
        description: "Não foi possível registrar a revisão. Tente novamente.",
        variant: "destructive",
      })
    } finally {
      setLoading(false)
    }
  }

  return (
    <Button onClick={handleReview} disabled={loading} className="gap-2">
      <CheckCircle className="h-4 w-4" />
      {loading ? "Registrando..." : "Marcar como Revisado"}
    </Button>
  )
}
