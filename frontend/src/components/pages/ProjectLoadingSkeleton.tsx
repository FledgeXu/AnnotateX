import { Skeleton } from "../ui/skeleton";


export const ProjectLoadingSkeleton = () => (
    <div className="space-y-2">
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-2/3" />
        <Skeleton className="h-4 w-1/3" />
    </div>
);
