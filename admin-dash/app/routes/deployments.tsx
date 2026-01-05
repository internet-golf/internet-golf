import { Typography } from "antd";
import type { Route } from "./+types/deployments";

export function meta({}: Route.MetaArgs) {
  return [{ title: "A Second Page" }, { name: "description", content: "Here is a Second Page!" }];
}

export default function Home() {
  return <Typography.Paragraph>This page is not the home page!</Typography.Paragraph>;
}
