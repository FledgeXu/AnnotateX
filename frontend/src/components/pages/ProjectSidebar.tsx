import { useQuery } from "@tanstack/react-query";
import { useStoreActions } from "easy-peasy";
import { Plus, SearchIcon } from "lucide-react";
import { useEffect } from "react";
import { toast } from "sonner";
import { CreateProjectDialog } from "@/components/pages/CreateProjectDialog";
import { SearchInput } from "@/components/pages/IconInput";
import { ProjectLoadingSkeleton } from "@/components/pages/ProjectLoadingSkeleton";
import { ProjectsFilterSelect } from "@/components/pages/ProjectsFilterSelect";
import { ProjectsList } from "@/components/pages/ProjectsList";
import { ProjectsSortSelect } from "@/components/pages/ProjectsSortSelect";
import { Button } from "@/components/ui/button";
import { createAPI } from "@/config";
import type { Project, Response } from "@/models";
import { store } from "@/store";
import type { StoreModel } from "@/store/types";

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
            <div className="grid grid-cols-2 gap-4">
                <ProjectsSortSelect className="w-full" />
                <ProjectsFilterSelect className="w-full" />
            </div>
            <div className="flex-1 min-h-0">
                {isSuccess && <ProjectsList />}
                {isFetching && <ProjectLoadingSkeleton />}
            </div>
        </div>
    );
};
