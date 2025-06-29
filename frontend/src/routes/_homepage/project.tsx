import { createFileRoute } from "@tanstack/react-router";
import { SearchIcon, SquarePlus } from "lucide-react";
import { SearchInput } from "@/components/pages/IconInput";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";

export const Route = createFileRoute("/_homepage/project")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="h-full flex gap-2">
      <div className="h-full w-xs flex flex-col">
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
        <div className="flex-1 min-h-0">
          <ScrollArea className="flex-1 h-full">
            <div className="p-2">
              {Array.from({ length: 200 }).map((_, i) => (
                <div key={i}>HHHH</div>
              ))}
            </div>
          </ScrollArea>
        </div>
      </div>
      <div className="bg-gray-100 rounded-sm border w-full h-full">HHH</div>
    </div>
  );
}
