import { useStoreState, useStoreActions } from "easy-peasy";
import * as React from "react";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import type { StoreModel } from "@/store/types";
type SearchInputProps = React.ComponentPropsWithoutRef<typeof Input> & {
    startIcon?: React.ReactNode;
};

export const SearchInput = React.forwardRef<HTMLInputElement, SearchInputProps>(
    ({ className, startIcon, ...props }, ref) => {
        const searchQuery = useStoreState<StoreModel>(
            (state) => state.projects.searchQuery,
        );
        const setSearchQuery = useStoreActions<StoreModel>(
            (actions) => actions.projects.setSearchQuery,
        );

        return (
            <div className={cn("flex items-center gap-2 relative", className)}>
                {startIcon && (
                    <span className="absolute left-3 text-gray-500 pointer-events-none">
                        {startIcon}
                    </span>
                )}
                <Input
                    ref={ref}
                    type="search"
                    className={cn("pl-9", startIcon ? "pl-9" : "")}
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    {...props}
                />
            </div>
        );
    },
);

SearchInput.displayName = "SearchInput";
