import { createFileRoute, Outlet } from "@tanstack/react-router";
import { ProjectSidebar } from "@/components/pages/ProjectSidebar";

export const Route = createFileRoute("/_homepage/project")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="h-full flex gap-2">
      <ProjectSidebar />
      <div className="rounded-sm border w-full h-full shadow-md p-3">
        <Outlet />
      </div>
    </div>
  );
}
