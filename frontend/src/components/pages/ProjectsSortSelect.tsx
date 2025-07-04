import { useStoreActions, useStoreState } from "easy-peasy";
import { ArrowUpDown } from "lucide-react";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";

import { PROJECT_SORT_MODES } from "@/models";
import type { StoreModel } from "@/store/types";

export const ProjectsSortSelect = ({ className }: { className?: string }) => {
    const sortMode = useStoreState<StoreModel>(
        (state) => state.projects.sortMode,
    );
    const setSortMode = useStoreActions<StoreModel>(
        (actions) => actions.projects.setSortMode,
    );

    return (
        <Select value={sortMode} onValueChange={(value) => setSortMode(value)}>
            <SelectTrigger className={className}>
                <ArrowUpDown />
                <SelectValue placeholder="Sort" />
            </SelectTrigger>
            <SelectContent>
                {PROJECT_SORT_MODES.map((mode) => (
                    <SelectItem value={mode} key={mode}>
                        {mode}
                    </SelectItem>
                ))}
            </SelectContent>
        </Select>
    );
};
