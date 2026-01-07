import { Typography } from "antd";
import type { Route } from "./+types/deployments";
import { redirect, useLoaderData } from "react-router";
import { getDeployments } from "~/api-calls/generated/golfComponents";

export function meta({}: Route.MetaArgs) {
  return [{ title: "A Second Page" }, { name: "description", content: "Here is a Second Page!" }];
}

export async function clientLoader() {
  const deployments = await getDeployments();
  console.log(deployments);
  return deployments;
}

export default function Home({ loaderData: { deployments } }: Route.ComponentProps) {
  const data = useLoaderData();
  console.log("in deployments:", data);
  return (
    <div>
      <Typography.Paragraph>This page is not the home page!</Typography.Paragraph>
      <code>{JSON.stringify(deployments, null, 4)}</code>
    </div>
  );
}
