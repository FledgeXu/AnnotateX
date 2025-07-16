import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { AxiosError } from "axios";
import { useStoreActions } from "easy-peasy";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import z from "zod";
import { Button } from "@/components/ui/button";
import {
    Dialog,
    DialogClose,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";
import {
    Form,
    FormControl,
    FormDescription,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { createAPI } from "@/config";
import type { Response, Project } from "@/models";
import { store } from "@/store";
import type { StoreModel } from "@/store/types";

const projectModality = ["2D", "3D", "audio", "text"] as const;

const schema = z.object({
    name: z.string().min(3, "Project is required"),
    description: z.string(),
    modality: z.enum(projectModality),
});
type FormData = z.infer<typeof schema>;

const useCreateProjectForm = () => {
    return useForm<FormData>({
        resolver: zodResolver(schema),
        defaultValues: {
            name: "",
            modality: "2D",
            description: "",
        },
    });
};

export const createProject = async (
    data: FormData,
): Promise<Response<Project>> => {
    const api = createAPI(store);
    const res = await api.post("/v1/projects/create", {
        name: data.name,
        modality: data.modality,
        description: data.description,
    });
    return res.data;
};

export const useCreateProjectMutation = (setOpen: (value: boolean) => void) => {
    const clearProjects = useStoreActions<StoreModel>(
        (state) => state.projects.clearProjects,
    );
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: createProject,
        onSuccess: (_) => {
            setOpen(false);
            queryClient.removeQueries({ queryKey: ["useFetchingProjects"] });
            clearProjects();
        },
        onError: (error: AxiosError<{ message?: string }>) => {
            const message = error?.response?.data?.message || "Failed to create.";
            toast.error(`Create failed: ${message}`);
        },
    });
};

export const CreateProjectDialog = ({
    children,
}: {
    children: React.ReactNode;
}) => {
    const [open, setOpen] = useState(false);
    const form = useCreateProjectForm();
    const createProjectMutation = useCreateProjectMutation(setOpen);
    const onSubmit = form.handleSubmit((data) =>
        createProjectMutation.mutate(data),
    );

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>{children}</DialogTrigger>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Create Project</DialogTitle>
                    <DialogDescription>
                        Create a new project by providing its name, description, and other
                        details.
                    </DialogDescription>
                </DialogHeader>
                <Form {...form}>
                    <form onSubmit={onSubmit} autoComplete="off">
                        <div className="flex flex-col gap-4">
                            <FormField
                                control={form.control}
                                name="name"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Project Name</FormLabel>
                                        <FormControl>
                                            <Input placeholder="Enter project name." {...field} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="modality"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Modality</FormLabel>
                                        <Select
                                            onValueChange={field.onChange}
                                            defaultValue={field.value}
                                        >
                                            <FormControl>
                                                <SelectTrigger>
                                                    <SelectValue placeholder="Select a verified email to display" />
                                                </SelectTrigger>
                                            </FormControl>
                                            <SelectContent>
                                                {projectModality.map((t, _) => (
                                                    <SelectItem key={t} value={t}>
                                                        {t}
                                                    </SelectItem>
                                                ))}
                                            </SelectContent>
                                        </Select>
                                        <FormDescription>The category of project.</FormDescription>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name="description"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Project Description</FormLabel>
                                        <FormControl>
                                            <Textarea
                                                placeholder="Enter project description."
                                                {...field}
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        </div>
                    </form>
                    <DialogFooter className="flex justify-end-safe">
                        <DialogClose asChild>
                            <Button type="button" variant="secondary">
                                Close
                            </Button>
                        </DialogClose>
                        <Button onClick={() => onSubmit()}>Create</Button>
                    </DialogFooter>
                </Form>
            </DialogContent>
        </Dialog>
    );
};
