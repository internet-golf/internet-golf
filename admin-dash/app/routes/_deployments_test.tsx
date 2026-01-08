import type { Route } from "./+types/deployments";
import { getDeployments, type GetDeploymentResponse } from "~/api-calls/generated/golfComponents";
import { DeploymentCard } from "~/components/DeploymentCard";

export function meta({}: Route.MetaArgs) {
  return [{ title: "Deployments" }];
}

export async function clientLoader() {
  return [
    {
      url: "internet-golf-test.local",
      tags: [],
      preserveExternalPath: false,
      meta: {
        title: "This is a Test",
        description: "An Internet Golf test page.",
        image: "https://mitch.website/_astro/circuit.D7sp5_r1_2msTYU.webp",
      },
      type: "StaticSite",
      createdAt: "0001-01-01T00:00:00Z",
      updatedAt: "2026-01-08T05:48:38Z",
      serverContentLocation:
        "C:/Users/Mitch/.internetgolf/internet-golf-test-local/accf561a85cb676c08a485e73dcf60fb",
      spaMode: false,
    },
    {
      url: "dash.internet-golf-test.local",
      tags: ["system"],
      preserveExternalPath: false,
      meta: { title: "", description: "", image: "" },
      type: "StaticSite",
      createdAt: "0001-01-01T00:00:00Z",
      updatedAt: "0001-01-01T00:00:00Z",
      serverContentLocation: "C:/Users/Mitch/.internetgolf/dashboard",
      spaMode: true,
    },
    {
      url: "fake.internet-golf-test.local",
      externalSource: "toBeOfUse/thingy",
      externalSourceType: "GithubRepo",
      tags: [],
      preserveExternalPath: false,
      meta: { title: "", description: "", image: "" },
      type: "Empty",
      createdAt: "0001-01-01T00:00:00Z",
      updatedAt: "2026-01-08T03:34:36Z",
      noContentYet: true,
    },
    {
      url: "other.internet-golf-test.local",
      tags: [],
      preserveExternalPath: false,
      meta: { title: "", description: "", image: "" },
      type: "Alias",
      createdAt: "0001-01-01T00:00:00Z",
      updatedAt: "0001-01-01T00:00:00Z",
      aliasedTo: "internet-golf-test.local",
      redirect: true,
    },
  ] satisfies GetDeploymentResponse[];
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
