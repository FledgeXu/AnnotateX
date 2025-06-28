import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { IndexNav } from "@/components/pages/IndexNav";

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
    <>
      <IndexNav />
      <Outlet />
    </>
  );
}
