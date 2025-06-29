import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { HomePageSidebar } from "@/components/pages/HomePageSidebar";
import { IndexNav } from "@/components/pages/IndexNav";
import { SidebarProvider } from "@/components/ui/sidebar";

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
      <HomePageSidebar />
      <div className="w-full">
        <IndexNav />
        <Outlet />
      </div>
    </SidebarProvider>
  );
}
