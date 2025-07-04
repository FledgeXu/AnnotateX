import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { IndexNav } from "@/components/pages/IndexNav";
import { SidebarProvider } from "@/components/ui/sidebar";
import { HomePageSidebar } from "@/components/pages/HomePageSidebar";

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
    <SidebarProvider className="h-full">
      <HomePageSidebar />
      <div className="w-full h-full flex flex-col">
        <IndexNav />
        <div className="flex-1 min-h-0 px-4 py-2">
          <Outlet />
        </div>
      </div>
    </SidebarProvider>
  );
}
