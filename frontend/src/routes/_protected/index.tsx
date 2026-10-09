import { getApiListFilesOptions } from "#/client/gen/@tanstack/react-query.gen.ts";
import type { FileItem } from "#/client/gen/index.ts";
import { Button } from "#/components/ui/button.tsx";
import { env } from "#/config/env.ts";
import { FileTable } from "#/features/explorer/components/file-table.tsx";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useEffect } from "react";

export const Route = createFileRoute("/_protected/")({
  component: Home,
});

const items: FileItem[] = [
  {
    name: "Test",
    is_dir: true,
    mod_time: new Date().toISOString(),
    path: "/",
    size: 0,
  },
  {
    name: "file.txt",
    is_dir: false,
    mod_time: new Date().toISOString(),
    path: "/",
    size: 100,
  },
];

function Home() {
  const { data } = useQuery(getApiListFilesOptions());
  useEffect(() => {
    if (data) console.log("get data");
    console.log(env.PUBLIC_APP_URL);
  }, [data]);

  function handleClick() {}
  return (
    <div>
      <h1 className="font-bold text-4xl">Welcome to TanStack Start</h1>
      <p className="mt-4 text-lg">
        Edit <code>src/routes/index.tsx</code> to get started.
      </p>

      <Button onClick={handleClick}>Click</Button>

      <FileTable data={items} />
    </div>
  );
}
