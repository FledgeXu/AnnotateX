import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import type { IUserAuth } from "@/user-auth";

interface IRouterContext {
  userAuth: IUserAuth;
}

export const Route = createRootRouteWithContext<IRouterContext>()({
  component: () => (
    <div className="h-dvh">
      <Outlet />
      {import.meta.env.DEV && <TanStackRouterDevtools />}
    </div>
  ),
});
