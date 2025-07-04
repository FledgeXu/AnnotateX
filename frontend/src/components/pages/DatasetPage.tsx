import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export const DatasetPage = () => {
    return (
        <div>
            <div className="flex justify-between">
                <Input placeholder="Dataset name" className="w-1/2" />
                <Button> Add </Button>
            </div>
        </div>
    );
};
