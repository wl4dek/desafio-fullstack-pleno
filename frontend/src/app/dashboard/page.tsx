import { SummaryCards } from "@/features/dashboard/components/SummaryCards"

export default function DashboardPage() {
  return (
    <div className="space-y-5">
      <div>
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <p className="text-sm text-neutral-500">Visão geral do acompanhamento</p>
      </div>
      <SummaryCards />
    </div>
  )
}
