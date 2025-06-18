import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_authenticated")({
  beforeLoad: async ({ context }) => {
    if (await context.userAuth.isLogin()) {
      throw redirect({ to: "/login" });
    }
  },
});
