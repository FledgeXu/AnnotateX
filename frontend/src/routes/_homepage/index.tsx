import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_homepage/')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/_homepage/"!</div>
}
