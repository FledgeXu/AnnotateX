import { createFileRoute } from "@tanstack/react-router";
import { Plus } from "lucide-react";
import { CreateProjectDialog } from "@/components/pages/CreateProjectDialog";
import { Button } from "@/components/ui/button";

export const Route = createFileRoute("/_homepage/project/")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="flex justify-center h-full">
      <div className="flex flex-col justify-center items-center gap-2">
        <div className="flex flex-col justify-center items-center">
          <div className="text-2xl font-bold">Please select a project</div>
          <div className="text-gray-500">
            Or create a new project to begin your annotation journey
          </div>
        </div>

        <CreateProjectDialog>
          <Button variant="outline">
            <Plus />
            New Project
          </Button>
        </CreateProjectDialog>
      </div>
    </div>
  );
}
