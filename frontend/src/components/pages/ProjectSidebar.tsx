import { useQuery } from "@tanstack/react-query";
import { SearchIcon, SquarePlus } from "lucide-react";
import { SearchInput } from "@/components/pages/IconInput";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { createAPI } from "@/config";
import type { ProjectResponse, Response } from "@/models";
import { store } from "@/store";
import { localizedDateFromISO } from "@/utils/date";
const statusColorMap = new Map<string, string>([
    ["active", "bg-green-600 text-white"],
    ["archive", "bg-gray-200 text-gray-800"],
]);

export const ProjectSidebar = () => {
    const { isPending, error, data, isFetching } = useQuery<
        Response<ProjectResponse>
    >({
        queryKey: ["queryProjects"],
        queryFn: async () => {
            const api = createAPI(store);
            const res = await api.get("/v1/projects/list");
            console.log(res.data);
            return res.data;
        },
    });

    return (
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
                    {data?.data.results.map((project) => (
                        <a href="/project/" key={project.id}>
                            <div className="p-2 w-full hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50 rounded-sm">
                                <div className="flex justify-between items-center">
                                    <div className="pb-2 font-medium">{project.name}</div>
                                    <Badge className={statusColorMap.get(project.status)}>
                                        {project.status}
                                    </Badge>
                                </div>
                                <span className="text-sm text-gray-500">{localizedDateFromISO(project.created_at)}</span>
                            </div>
                        </a>
                    ))}
                </ScrollArea>
            </div>
        </div>
    );
};
