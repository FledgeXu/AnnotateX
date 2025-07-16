import { useInfiniteQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { useStoreActions } from "easy-peasy";
import { useStoreState } from "easy-peasy";
import { useEffect } from "react";
import { useInView } from "react-intersection-observer";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { createAPI } from "@/config";
import type { Response, Paginated, Project } from "@/models";
import { store } from "@/store";
import type { StoreModel } from "@/store/types";
import { localizedDateFromISO } from "@/utils/date";
const statusColorMap = new Map<string, string>([
    ["active", "bg-green-600 text-white"],
    ["archive", "bg-gray-200 text-gray-800"],
]);

const LIMIT = 10;

const useFetchingProjects = () => {
    const updateProjects = useStoreActions<StoreModel>(
        (state) => state.projects.updateProjects,
    );

    const { data, fetchNextPage, hasNextPage, isFetchingNextPage } =
        useInfiniteQuery<Response<Paginated<Project>>>({
            queryKey: ["useFetchingProjects"],
            queryFn: async ({ pageParam }: { pageParam: unknown }) => {
                const offset = typeof pageParam === "number" ? pageParam : 0;
                const api = createAPI(store);
                const res = await api.get("/v1/projects/list", {
                    params: { offset, limit: LIMIT },
                });
                return res.data;
            },
            initialPageParam: 0,
            getNextPageParam: (lastPage) => {
                const { offset, results = [] } = lastPage.data;
                const nextOffset = offset + results.length;
                return results.length === 0 ? undefined : nextOffset;
            },
        });

    useEffect(() => {
        if (data != undefined) {
            updateProjects(data.pages.flatMap((page) => page.data.results));
        }
    }, [data, updateProjects]);

    const { ref, inView } = useInView();

    useEffect(() => {
        if (inView && hasNextPage && !isFetchingNextPage) {
            fetchNextPage();
        }
    }, [inView, hasNextPage, fetchNextPage, isFetchingNextPage]);

    return { hasNextPage, isFetchingNextPage, ref };
};

export const ProjectsList = () => {
    const projects = useStoreState<StoreModel>(
        (state) => state.projects.visibleProjects,
    ) as Project[];

    const { isFetchingNextPage, hasNextPage, ref } = useFetchingProjects();

    return (
        <ScrollArea className="h-full pr-3">
            {projects.map((project, index) => {
                const isLoadingElement = index === Math.max(projects.length - 5, 0);
                return (
                    <Link
                        to="/project/$id"
                        params={{ id: String(project.id) }}
                        key={project.id}
                    >
                        <div
                            className="w-full hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50 rounded-sm"
                            ref={isLoadingElement ? ref : undefined}
                        >
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
                );
            })}
            <div ref={ref} className="flex justify-center">
                {isFetchingNextPage && hasNextPage && (
                    <Skeleton className="h-4 w-3/4" />
                )}
                {!hasNextPage && (
                    <div className="text-gray-500 text-sm">No more data</div>
                )}
            </div>
        </ScrollArea>
    );
};
