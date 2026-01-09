import { Typography } from "antd";
import { type GetDeploymentResponse } from "~/api-calls/generated/golfComponents";
import { DeploymentCard } from "~/components/DeploymentCard";

export default function DeploymentTest() {
  const deployments = [
    // static site with successfully collected metadata and repo:
    {
      name: "",
      url: "internet-golf-test.local",
      tags: [],
      preserveExternalPath: false,
      externalSource: "toBeOfUse/cool-website",
      externalSourceType: "Github",
      type: "StaticSite",
      createdAt: "0001-01-01T00:00:00Z",
      updatedAt: "2026-01-08T05:48:38Z",
      meta: {
        title: "Site with Metadata and Repo",
        description: "This is the <meta> data for this site.",
        image: "https://mitch.website/_astro/circuit.D7sp5_r1_2msTYU.webp",
      },
      serverContentLocation:
        "C:/Users/Mitch/.internetgolf/internet-golf-test-local/accf561a85cb676c08a485e73dcf60fb",
      spaMode: false,
    },
    // static site with successfully collected metadata except for the image:
    {
      name: "",
      url: "internet-golf-test.local",
      tags: [],
      preserveExternalPath: false,
      type: "StaticSite",
      createdAt: "0001-01-01T00:00:00Z",
      updatedAt: "2026-01-08T05:48:38Z",
      meta: {
        title: "Site with Metadata Minus Image",
        description: "This is the <meta> data for this site.",
        image: "",
      },
      serverContentLocation:
        "C:/Users/Mitch/.internetgolf/internet-golf-test-local/accf561a85cb676c08a485e73dcf60fb",
      spaMode: false,
    },
    // static site with successfully collected metadata except for the title and
    // description, and a name in the deployments db:
    {
      name: "Site With Name In DB",
      url: "internet-golf-test.local",
      tags: [],
      preserveExternalPath: false,
      type: "StaticSite",
      createdAt: "0001-01-01T00:00:00Z",
      updatedAt: "2026-01-08T05:48:38Z",
      meta: {
        title: "",
        description: "",
        image: "https://mitch.website/_astro/circuit.D7sp5_r1_2msTYU.webp",
      },
      serverContentLocation:
        "C:/Users/Mitch/.internetgolf/internet-golf-test-local/accf561a85cb676c08a485e73dcf60fb",
      spaMode: false,
    },
    // static site with no metadata:
    {
      name: "",
      url: "internet-golf-test.local",
      tags: [],
      preserveExternalPath: false,
      type: "StaticSite",
      createdAt: "0001-01-01T00:00:00Z",
      updatedAt: "2026-01-08T05:48:38Z",
      meta: {
        title: "",
        description: "",
        image: "",
      },
      serverContentLocation:
        "C:/Users/Mitch/.internetgolf/internet-golf-test-local/accf561a85cb676c08a485e73dcf60fb",
      spaMode: false,
    },
    // empty deployment:
    {
      name: "",
      url: "empty.internet-golf-test.local",
      tags: [],
      preserveExternalPath: false,
      type: "Empty",
      createdAt: "0001-01-01T00:00:00Z",
      updatedAt: "2026-01-08T21:49:12Z",
      meta: {
        title: "",
        description: "",
        image: "",
      },
      noContentYet: true,
    },
    // empty deployment with name and github repo:
    {
      name: "Empty deployment w/ name & repo",
      url: "fake.internet-golf-test.local",
      externalSource: "toBeOfUse/empty",
      externalSourceType: "Github",
      tags: [],
      preserveExternalPath: false,
      type: "Empty",
      createdAt: "0001-01-01T00:00:00Z",
      updatedAt: "2026-01-08T21:49:12Z",
      meta: {
        title: "",
        description: "",
        image: "",
      },
      noContentYet: true,
    },
    // alias deployment:
    {
      name: "",
      url: "alias.internet-golf-test.local",
      tags: [],
      preserveExternalPath: false,
      type: "Alias",
      createdAt: "0001-01-01T00:00:00Z",
      updatedAt: "0001-01-01T00:00:00Z",
      meta: {
        title: "",
        description: "",
        image: "",
      },
      aliasedTo: "internet-golf-test.local",
      redirect: false,
    },
    // redirect deployment:
    {
      name: "Redirect to main site",
      url: "alias.internet-golf-test.local",
      tags: [],
      preserveExternalPath: false,
      type: "Alias",
      createdAt: "0001-01-01T00:00:00Z",
      updatedAt: "0001-01-01T00:00:00Z",
      meta: {
        title: "",
        description: "",
        image: "",
      },
      aliasedTo: "internet-golf-test.local",
      redirect: true,
    },
  ] satisfies GetDeploymentResponse[];

  return !!deployments?.length ? (
    <>
      <Typography.Title level={2}>Deployments</Typography.Title>
      <div className="grid gap-8 w-full grid-cols-1 md:grid-cols-2">
        {deployments?.map((d) => (
          <DeploymentCard key={d.url} deployment={d} />
        ))}
      </div>
    </>
  ) : (
    // TODO: nicer content for this
    <div>You have no deployments yet!</div>
  );
}
