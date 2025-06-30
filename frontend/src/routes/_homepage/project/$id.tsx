import { createFileRoute, useParams } from "@tanstack/react-router";

export const Route = createFileRoute("/_homepage/project/$id")({
  component: RouteComponent,
});

function RouteComponent() {
  const { id } = useParams({ strict: false });
  return <div>Hello `/_homepage/project/${id}`!</div>;
}
