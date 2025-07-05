import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useParams } from "@tanstack/react-router";
import { DatasetPage } from "@/components/pages/DatasetPage";
import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { createAPI } from "@/config";
import type { Project } from "@/models";
import { store } from "@/store";

export const Route = createFileRoute("/_homepage/project/$id")({
  component: RouteComponent,
});

function RouteComponent() {
  const { id } = useParams({ strict: false });
  const { data: project } = useQuery<Project>({
    queryKey: ["fetchProjectInfo", id],
    queryFn: async () => {
      const api = createAPI(store);
      const res = await api.get(`/v1/projects/${id}`);
      return res.data.data;
    },
  });

  return (
    <Tabs defaultValue="account">
      <div className="flex justify-between items-center">
        <span className="text-2xl font-semibold">{project?.name ?? ""}</span>
        <TabsList>
          <TabsTrigger value="batches">Batches</TabsTrigger>
          <TabsTrigger value="dataset">Dataset</TabsTrigger>
        </TabsList>
      </div>
      <Separator />
      <TabsContent value="batches">
        Make changes to your account here.
      </TabsContent>
      <TabsContent value="dataset">
        <DatasetPage />
      </TabsContent>
    </Tabs>
  );
}
