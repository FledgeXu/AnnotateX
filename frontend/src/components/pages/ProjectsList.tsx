import { Link } from "@tanstack/react-router";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";

import { Separator } from "@/components/ui/separator";
import type { Project } from "@/models";
import { localizedDateFromISO } from "@/utils/date";
const statusColorMap = new Map<string, string>([
    ["active", "bg-green-600 text-white"],
    ["archive", "bg-gray-200 text-gray-800"],
]);

type ProjectListProps = {
    projects: Project[];
};
export const ProjectsList = ({ projects }: ProjectListProps) => (
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
