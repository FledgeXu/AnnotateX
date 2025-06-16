import { StrictMode } from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider, createRouter } from "@tanstack/react-router";
import { Theme } from "@radix-ui/themes";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import "@/main.css";
import "@radix-ui/themes/styles.css";

// Import the generated route tree
import { routeTree } from "@/routeTree.gen";
import { StoreProvider } from "easy-peasy";
import { store } from "@/store";

// Create a new router instance
const router = createRouter({ routeTree });

// Register the router instance for type safety
declare module "@tanstack/react-router" {
    interface Register {
        router: typeof router;
    }
}

const queryClient = new QueryClient();

// Render the app
const rootElement = document.getElementById("root")!;
if (!rootElement.innerHTML) {
    const root = ReactDOM.createRoot(rootElement);
    root.render(
        <StrictMode>
            <StoreProvider store={store}>
                <QueryClientProvider client={queryClient}>
                    <Theme>
                        <RouterProvider router={router} />
                    </Theme>
                    {import.meta.env.DEV && <ReactQueryDevtools />}
                </QueryClientProvider>
            </StoreProvider>
        </StrictMode>,
    );
}
