import { createFileRoute } from "@tanstack/react-router";
import { SearchIcon, SquarePlus } from "lucide-react";
import { SearchInput } from "@/components/pages/IconInput";
import { Button } from "@/components/ui/button";

export const Route = createFileRoute("/_homepage/project")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="flex h-full w-full gap-2">
      <div className="w-xs">
        <div className="flex justify-between">
          <h1 className="text-2xl font-bold">Project</h1>
          <Button variant={"ghost"}>
            <SquarePlus />
          </Button>
        </div>
        <SearchInput
          placeholder="Search..."
          startIcon={<SearchIcon className="w-4 h-4" />}
        />
      </div>

      <div className="bg-amber-300 w-full h-full">HHH</div>
    </div>
  );
}
