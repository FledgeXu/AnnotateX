import { createFileRoute } from "@tanstack/react-router";
import { ProjectSidebar } from "@/components/ProjectSidebar";

export const Route = createFileRoute("/_homepage/project")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="h-full flex gap-2">
      <ProjectSidebar />
      <div className="bg-gray-100 rounded-sm border w-full h-full">HHH</div>
    </div>
  );
}
