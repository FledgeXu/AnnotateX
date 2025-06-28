import { createFileRoute } from "@tanstack/react-router";

import { IndexNav } from "@/components/pages/IndexNav";

export const Route = createFileRoute("/_authenticated/")({
  component: Index,
});

function Index() {
  return <IndexNav />;
}
