import { ChildDetailPage } from "./child-detail-page"

export default async function ChildPage({
  params,
}: {
  params: Promise<{ id: string }>
}) {
  const { id } = await params
  return <ChildDetailPage id={id} />
}
