import { createFileRoute, useParams } from "@tanstack/react-router";
import { DatasetPage } from "@/components/pages/DatasetPage";
import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

export const Route = createFileRoute("/_homepage/project/$id")({
  component: RouteComponent,
});

function RouteComponent() {
  const { id } = useParams({ strict: false });
  return (
    <Tabs defaultValue="account">
      <div className="flex justify-between items-center">
        <span className="text-2xl font-semibold">{id}</span>
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
