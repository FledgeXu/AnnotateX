import { useQueryClient } from "@tanstack/react-query";
import { useStoreActions, useStoreState } from "easy-peasy";
import { ArrowUpDown } from "lucide-react";
import { useTranslation } from "react-i18next";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";

import { PROJECT_SORT_MODES, type ProjectSortMode } from "@/models";
import type { StoreModel } from "@/store/types";

export const ProjectsSortSelect = ({ className }: { className?: string }) => {
    const sortMode = useStoreState<StoreModel>(
        (state) => state.projects.sortMode,
    );
    const setSortMode = useStoreActions<StoreModel>(
        (actions) => actions.projects.setSortMode,
    );
    const { t } = useTranslation();

    const clearProjects = useStoreActions<StoreModel>(
        (state) => state.projects.clearProjects,
    );
    const queryClient = useQueryClient();

    const onValueChange = (value: ProjectSortMode) => {
        queryClient.removeQueries({ queryKey: ["useFetchingProjects"] });
        clearProjects();
        setSortMode(value);
    };

    return (
        <Select value={sortMode} onValueChange={onValueChange}>
            <SelectTrigger className={className}>
                <ArrowUpDown />
                <SelectValue placeholder="Sort" />
            </SelectTrigger>
            <SelectContent>
                {PROJECT_SORT_MODES.map((mode) => (
                    <SelectItem value={mode} key={mode}>
                        {t(mode)}
                    </SelectItem>
                ))}
            </SelectContent>
        </Select>
    );
};
