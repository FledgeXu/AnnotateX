import { createFileRoute, useParams } from "@tanstack/react-router";
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
          <TabsTrigger value="account">Batches</TabsTrigger>
          <TabsTrigger value="password">Dataset</TabsTrigger>
        </TabsList>
      </div>
      <TabsContent value="account">
        Make changes to your account here.
      </TabsContent>
      <TabsContent value="password">Change your password here.</TabsContent>
    </Tabs>
  );
}
