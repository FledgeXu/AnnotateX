import { ListFilter } from "lucide-react";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";

export const ProjectsFilterSelect = ({ className }: { className: string }) => {
    return (
        <Select>
            <SelectTrigger className={className}>
                <ListFilter />
                <SelectValue placeholder="Filter" />
            </SelectTrigger>
            <SelectContent>
                <SelectItem value="light">Name</SelectItem>
                <SelectItem value="dark">Time</SelectItem>
            </SelectContent>
        </Select>
    );
};
