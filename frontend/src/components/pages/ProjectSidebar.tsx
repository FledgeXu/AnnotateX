import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { Plus, SearchIcon } from "lucide-react";
import { SearchInput } from "@/components/pages/IconInput";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { createAPI } from "@/config";
import type { Project, Response } from "@/models";
import { store } from "@/store";
import { localizedDateFromISO } from "@/utils/date";
const statusColorMap = new Map<string, string>([
    ["active", "bg-green-600 text-white"],
    ["archive", "bg-gray-200 text-gray-800"],
]);

type ProjectListProps = {
    projects: Project[];
};
const ProjectList = ({ projects }: ProjectListProps) => (
    <ScrollArea className="h-full">
        {projects.map((project, index) => (
            <Link
                to="/project/$id"
                params={{ id: String(project.id) }}
                key={project.id}
            >
                <div className="p-2 w-full hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50 rounded-sm">
                    <div className="flex justify-between items-center">
                        <div className="pb-2 font-medium">{project.name}</div>
                        <Badge className={statusColorMap.get(project.status)}>
                            {project.status}
                        </Badge>
                    </div>
                    <span className="text-sm text-gray-500">
                        {localizedDateFromISO(project.created_at)}
                    </span>
                </div>
                {index < projects.length - 1 && <Separator className="m-2" />}
            </Link>
        ))}
    </ScrollArea>
);

export const ProjectSidebar = () => {
    const { isPending, error, data, isFetching } = useQuery<Response<Project[]>>({
        queryKey: ["queryProjects"],
        queryFn: async () => {
            const api = createAPI(store);
            const res = await api.get("/v1/projects/list");
            return res.data;
        },
    });

    return (
        <div className="h-full w-sm flex flex-col gap-2">
            <div className="flex justify-between">
                <h1 className="text-2xl font-bold">Project</h1>
                <Button variant={"outline"}>
                    <Plus />
                </Button>
            </div>
            <SearchInput
                placeholder="Search..."
                startIcon={<SearchIcon className="w-4 h-4" />}
            />
            <div className="flex-1 min-h-0">
                {data && <ProjectList projects={data.data} />}
            </div>
        </div>
    );
};
