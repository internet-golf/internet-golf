import type { Route } from "./+types/deployments";

export function meta({}: Route.MetaArgs) {
  return [{ title: "A Second Page" }, { name: "description", content: "Here is a Second Page!" }];
}

export default function Home() {
  return <p>This page is not the home page!</p>;
}
