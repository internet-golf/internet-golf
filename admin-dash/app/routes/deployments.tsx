import type { Route } from "./+types/deployments";
import { getDeployments } from "~/api-calls/generated/golfComponents";
import { DeploymentCard } from "~/components/DeploymentCard";

export function meta({}: Route.MetaArgs) {
  return [{ title: "Deployments" }];
}

export async function clientLoader() {
  const deployments = await getDeployments();
  return deployments;
}

export default function Home({ loaderData: { deployments } }: Route.ComponentProps) {
  return (
    <div className="grid gap-4 w-full grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
      {deployments?.map((d) => (
        <DeploymentCard key={d.url} deployment={d} />
      ))}
    </div>
  );
}
