import { ArrowUpDown } from "lucide-react";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";

export const ProjectsSortSelect = () => {
    return (
        <Select>
            <SelectTrigger className="w-full">
                <ArrowUpDown />
                <SelectValue placeholder="Sort" />
            </SelectTrigger>
            <SelectContent>
                <SelectItem value="light">Increase</SelectItem>
                <SelectItem value="dark">Decrease</SelectItem>
            </SelectContent>
        </Select>
    );
};
