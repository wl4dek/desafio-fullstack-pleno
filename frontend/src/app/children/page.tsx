import { ChildrenList } from "./children-list"

export default function ChildrenPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Crianças</h1>
        <p className="text-sm text-neutral-500">Listagem e filtros</p>
      </div>
      <ChildrenList />
    </div>
  )
}
