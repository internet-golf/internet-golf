import { Typography } from "antd";
import type { Route } from "./+types/deployments";
import { deleteDeployment, getDeployments } from "~/api-calls/generated/golfComponents";
import { DeploymentCard } from "~/components/DeploymentCard";
import { useFetcher } from "react-router";

export function meta({}: Route.MetaArgs) {
  return [{ title: "Deployments" }];
}

export async function clientLoader() {
  const deployments = await getDeployments();
  return deployments;
}

export async function clientAction({ request }: Route.ClientActionArgs) {
  const formData = await request.formData();

  if (request.method === "DELETE") {
    const url = formData.get("url") as string;
    await deleteDeployment({ pathParams: { url } });
  }
}

export default function Home({ loaderData: { deployments } }: Route.ComponentProps) {
  const fetcher = useFetcher();

  return !!deployments?.length ? (
    <>
      <Typography.Title level={2}>Deployments</Typography.Title>
      <div className="grid gap-8 w-full grid-cols-1 md:grid-cols-2">
        {deployments?.map((deployment) => (
          <DeploymentCard
            key={deployment.url}
            deployment={deployment}
            onDelete={() => {
              return fetcher.submit({ url: deployment.url }, { method: "DELETE" });
            }}
          />
        ))}
      </div>
    </>
  ) : (
    // TODO: nicer content for this
    <div>You have no deployments yet!</div>
  );
}
