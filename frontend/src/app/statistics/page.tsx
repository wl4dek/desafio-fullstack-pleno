import { AlertsArea } from "@/features/statistics/components/AlertsArea";

export default function StatisticsPage() {
  return (
    <div className="space-y-5">
      <div>
        <h1 className="text-2xl font-bold">Estatísticas</h1>
      </div>

      <AlertsArea />
    </div>
  )
}
