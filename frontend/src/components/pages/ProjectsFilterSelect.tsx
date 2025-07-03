import { ListFilter } from "lucide-react";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";

export const ProjectsFilterSelect = () => {
    return (
        <Select>
            <SelectTrigger className="w-full">
                <ListFilter />
                <SelectValue placeholder="Filter" />
            </SelectTrigger>
            <SelectContent>
                <SelectItem value="light">Name</SelectItem>
                <SelectItem value="dark">Create Time</SelectItem>
            </SelectContent>
        </Select>
    );
};
