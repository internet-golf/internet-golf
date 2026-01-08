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
    <div className="grid gap-8 w-full grid-cols-1 md:grid-cols-2">
      {deployments?.map((d) => (
        <DeploymentCard key={d.url} deployment={d} />
      ))}
    </div>
  );
}
