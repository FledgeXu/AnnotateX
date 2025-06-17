import type { IUserAuth } from "@/user-auth";
import {
    createRootRouteWithContext,
    Outlet,
} from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";

interface IRouterContext {
    userAuth: IUserAuth;
}

export const Route = createRootRouteWithContext<IRouterContext>()({
    component: () => (
        <>
            <Outlet />
            {import.meta.env.DEV && <TanStackRouterDevtools />}
        </>
    ),
});
