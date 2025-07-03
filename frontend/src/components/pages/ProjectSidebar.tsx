import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { useStoreActions, useStoreState } from "easy-peasy";
import { Plus, SearchIcon } from "lucide-react";
import { useEffect } from "react";
import { toast } from "sonner";
import { Skeleton } from "../ui/skeleton";
import { CreateProjectDialog } from "./CreateProjectDialog";
import { SearchInput } from "@/components/pages/IconInput";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { createAPI } from "@/config";
import type { Project, Response } from "@/models";
import { store } from "@/store";
import type { StoreModel } from "@/store/types";
import { localizedDateFromISO } from "@/utils/date";
const statusColorMap = new Map<string, string>([
    ["active", "bg-green-600 text-white"],
    ["archive", "bg-gray-200 text-gray-800"],
]);

type ProjectListProps = {
    projects: Project[];
};
const ProjectList = ({ projects }: ProjectListProps) => (
    <ScrollArea className="h-full pr-3">
        {projects.map((project, index) => (
            <Link
                to="/project/$id"
                params={{ id: String(project.id) }}
                key={project.id}
            >
                <div className="w-full hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50 rounded-sm">
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

const ProjectSkeleton = () => (
    <div className="space-y-2">
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-2/3" />
        <Skeleton className="h-4 w-1/3" />
    </div>
);

const useFetchingProjects = () => {
    const updateProjects = useStoreActions<StoreModel>(
        (state) => state.projects.updateProjects,
    );

    const { error, data, isFetching, isSuccess } = useQuery<Response<Project[]>>({
        queryKey: ["queryProjects"],
        queryFn: async () => {
            const api = createAPI(store);
            const res = await api.get("/v1/projects/list");
            return res.data;
        },
    });

    useEffect(() => {
        if (isSuccess && data) {
            updateProjects(data.data);
        }
    }, [isSuccess]);

    useEffect(() => {
        if (error) {
            toast.error("Fail to load projects.");
        }
    }, [error]);

    return [isFetching, isSuccess];
};

export const ProjectSidebar = () => {
    const [isFetching, isSuccess] = useFetchingProjects();

    const projects = useStoreState<StoreModel>(
        (state) => state.projects.projects,
    );

    return (
        <div className="h-full w-sm flex flex-col gap-2">
            <div className="flex justify-between">
                <h1 className="text-2xl font-bold">Project</h1>
                <CreateProjectDialog>
                    <Button variant={"outline"}>
                        <Plus />
                    </Button>
                </CreateProjectDialog>
            </div>
            <SearchInput
                placeholder="Search..."
                startIcon={<SearchIcon className="w-4 h-4" />}
            />
            <div className="flex-1 min-h-0">
                {isSuccess && <ProjectList projects={projects} />}
                {isFetching && <ProjectSkeleton />}
            </div>
        </div>
    );
};
