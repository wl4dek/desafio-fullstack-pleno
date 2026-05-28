"use client"

import { ReviewButton } from "./ReviewButton"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"
import { Button } from "@/components/ui/button"
import { ArrowLeft, RefreshCw } from "lucide-react"
import { useRouter } from "next/navigation"
import { useChild, useChildAreas } from "@/hooks/useChildren"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Education, Health, SocialAssistance } from "@/types"
import { formatDateBR } from "@/lib/utils"
import { Alerts, AlertType } from "@/stores/alert"

interface ChildDetailProps {
  id: string
}

export function ChildDetail({ id }: ChildDetailProps) {
  const { child, isLoading, isError, error } = useChild(id)
  const { areas, isLoading: isAreasLoading, isError: isAreasError, error: areasError } = useChildAreas(id)
  const router = useRouter()

  if (isLoading && isAreasLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-32 w-full" />
      </div>
    )
  }

  if (isError) {
    return (
      <Card className="p-6 text-center">
        <p className="text-red-500 mb-2">
          {error instanceof Error ? error.message : "Erro ao carregar dados"}
        </p>
        <Button variant="outline" size="sm">
          <RefreshCw className="h-4 w-4 mr-2" />
          Tentar novamente
        </Button>
      </Card>
    )
  }

  if (!child) {
    return (
      <Card className="p-6 text-center">
        <p className="text-neutral-500">Criança não encontrada</p>
        <Button variant="outline" size="sm" className="mt-2" onClick={() => router.push("/children")}>
          Voltar
        </Button>
      </Card>
    )
  }

  const areaData = {
    "Saúde": areas ? { ...areas.health, lastConsultation: formatDateBR(areas.health.lastConsultation), type: "health" } : {} as Health,
    "Educação": areas ? { ...areas.education, type: "education" } : {} as Education,
    "Assistência Social": areas ? { ...areas.socialAssistance, type: "social" } : {} as SocialAssistance,
  } as const
  const areasTab = Object.keys(areaData) as Array<keyof typeof areaData>
  console.log("Child areas data:", areas)

  return (
    <div className="space-y-6">
      <Button variant="ghost" size="sm" onClick={() => router.push("/children")} className="gap-2">
        <ArrowLeft className="h-4 w-4" />
        Voltar
      </Button>
      <div className="flex gap-4">
        <Card className="flex-1">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-xl">{child.name}</CardTitle>
                <p className="text-sm text-neutral-500 mt-1">
                  {child.neighborhood}
                </p>
              </div>
              <div className="flex items-center gap-2">
                {child.has_alert && <Badge variant="destructive">Alerta</Badge>}
                {child.reviewed ? (
                  <Badge variant="success">Revisado</Badge>
                ) : (
                  <Badge variant="warning">Pendente</Badge>
                )}
              </div>
            </div>
          </CardHeader>
          <CardContent className="text-sm text-neutral-600 space-y-1">
            <p>Idade: {child.age} anos</p>
            {child.reviewed_at && (
              <p>Última revisão: {new Date(child.reviewed_at).toLocaleDateString("pt-BR")}</p>
            )}
          </CardContent>
        </Card>

        <Tabs className="flex-1" defaultValue={areasTab[0]}>
          <TabsList className="w-full overflow-x-auto">
            {areasTab.map((area) => (
              <TabsTrigger key={area} value={area} className="flex-1">
                {area}
              </TabsTrigger>
            ))}
          </TabsList>
          {areasTab.map((area) => (
            <TabsContent key={area} value={area}>
              <AreaSection area={area} notes={areaData[area]} />
            </TabsContent>
          ))}
        </Tabs>
      </div>
      {!child.reviewed && (
        <Card>
          <CardContent className="pt-6">
            <ReviewButton childId={child.id} />
          </CardContent>
        </Card>
      )}
    </div>
  )
}

function AreaSection({ area, notes }: { area: string; notes: Health | SocialAssistance | Education }) {
  if (!notes) {
    return (
      <Card>
        <CardContent className="py-8 text-center">
          <div className="text-neutral-400 text-4xl mb-2">-</div>
          <p className="text-sm text-neutral-500">
            Sem acompanhamento registrado em {area}
          </p>
        </CardContent>
      </Card>
    )
  }

  // education
  if (notes.type === "education") {
    const educationNotes = notes as Education
    return (
      <Card>
        <CardContent className="flex flex-col gap-2 py-4">
          <div className="flex gap-2 text-sm whitespace-pre-wrap">Alertas: {educationNotes.alerts?.length > 0 ? educationNotes.alerts.map((alert, i) => (
            <Badge key={i} variant="destructive">{Alerts[alert as AlertType]}</Badge>
          )) : <span className="text-muted-foreground">Nenhum alerta</span>}</div>
          <p className="text-sm whitespace-pre-wrap">Escola: {educationNotes.schoolName}</p>
          <p className="text-sm whitespace-pre-wrap">Frequencia: {educationNotes.frequenciaPercent}</p>
        </CardContent>
      </Card>
    )
  }

  // health
  if (notes.type === "health") {
    const healthNotes = notes as Health
    return (
      <Card>
        <CardContent className="flex flex-col gap-2 py-4">
          <div className="flex gap-2 text-sm whitespace-pre-wrap">Alertas: {healthNotes.alerts?.length > 0 ? healthNotes.alerts.map((alert, i) => (
            <Badge key={i} variant="destructive">{Alerts[alert as AlertType]}</Badge>
          )) : <span className="text-muted-foreground">Nenhum alerta</span>}</div>
          <p className="text-sm whitespace-pre-wrap">Vacinado: {healthNotes.vaccinationsUpToDate ? (<Badge className="bg-green-50 text-green-700 dark:bg-green-950 dark:text-green-300">
            Sim
          </Badge>) : (<Badge className="bg-red-50 text-red-700 dark:bg-red-950 dark:text-red-300">
            Não
          </Badge>)}</p>
          <p className="text-sm whitespace-pre-wrap">Última consulta: {healthNotes.lastConsultation}</p>
        </CardContent>
      </Card>
    )
  }

  // social assistance
  const socialNotes = notes as SocialAssistance
  return (
    <Card>
      <CardContent className="flex flex-col gap-2 py-4">
        <div className="flex gap-2 text-sm whitespace-pre-wrap">Alertas: {socialNotes.alerts?.length > 0 ? socialNotes.alerts.map((alert, i) => (
          <Badge key={i} variant="destructive">{Alerts[alert as AlertType]}</Badge>
        )) : <span className="text-muted-foreground">Nenhum alerta</span>}</div>
        <p className="text-sm whitespace-pre-wrap">Benefício Ativo: {socialNotes.activeBenefit ? (<Badge className="bg-green-50 text-green-700 dark:bg-green-950 dark:text-green-300">
          Sim
        </Badge>) : (<Badge className="bg-red-50 text-red-700 dark:bg-red-950 dark:text-red-300">
          Não
        </Badge>)}</p>
        <p className="text-sm whitespace-pre-wrap">Cartão Único: {socialNotes.cadUnico ? (<Badge className="bg-green-50 text-green-700 dark:bg-green-950 dark:text-green-300">
          Sim
        </Badge>) : (<Badge className="bg-red-50 text-red-700 dark:bg-red-950 dark:text-red-300">
          Não
        </Badge>)}</p>
      </CardContent>
    </Card>
  )
}

