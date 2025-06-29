import { SearchIcon, SquarePlus } from "lucide-react";
import { SearchInput } from "@/components/pages/IconInput";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";

export const ProjectSidebar = () => {
    return (
        <>
            <div className="h-full w-sm flex flex-col gap-2">
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
        </>
    );
};
