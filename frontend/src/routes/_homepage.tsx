import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { IndexNav } from "@/components/pages/IndexNav";
import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarProvider,
} from "@/components/ui/sidebar";

export const Route = createFileRoute("/_homepage")({
  beforeLoad: async ({ context }) => {
    if (!(await context.userAuth.isLogin())) {
      throw redirect({ to: "/login" });
    }
  },
  component: AuthenticatedLayout,
});

function AuthenticatedLayout() {
  return (
    <SidebarProvider>
      <Sidebar>
        <SidebarHeader>AnnotateX</SidebarHeader>
        <SidebarContent />
      </Sidebar>
      <div className="w-full">
        <IndexNav />
        <Outlet />
      </div>
    </SidebarProvider>
  );
}
