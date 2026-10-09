import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/explorer/$sourceId/")({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/_protected/explorer/$storageId/"!</div>;
}
