import { getApiOptions } from "#/client/gen/@tanstack/react-query.gen.ts";
import { env } from "#/config/env.ts";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useEffect } from "react";

export const Route = createFileRoute("/")({
  component: Home,
});

function Home() {
  const { data } = useQuery(getApiOptions());
  useEffect(() => {
    if (data) console.log("get data");
    console.log(env.PUBLIC_APP_URL);
  }, [data]);

  function handleClick() {}
  return (
    <div className="p-8">
      <h1 className="text-4xl font-bold">Welcome to TanStack Start</h1>
      <p className="mt-4 text-lg">
        Edit <code>src/routes/index.tsx</code> to get started.
      </p>

      <button onClick={handleClick}></button>
    </div>
  );
}
