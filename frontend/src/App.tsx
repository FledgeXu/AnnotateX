import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { RouterProvider, createRouter } from "@tanstack/react-router";
import { Toaster } from "@/components/ui/sonner";
import "@/main.css";
import "@radix-ui/themes/styles.css";

// Import the generated route tree
import { routeTree } from "@/routeTree.gen";
import { StoreProvider } from "easy-peasy";
import { store } from "@/store";
import { useUserAuth } from "@/user-auth";

// Create a new router instance
const router = createRouter({
  routeTree,
  context: {
    userAuth: undefined!,
  },
});

// Register the router instance for type safety
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

const queryClient = new QueryClient();

export const App = () => {
  const userAuth = useUserAuth();
  return (
    <StoreProvider store={store}>
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} context={{ userAuth }} />
        <Toaster />
        {import.meta.env.DEV && <ReactQueryDevtools />}
      </QueryClientProvider>
    </StoreProvider>
  );
};
