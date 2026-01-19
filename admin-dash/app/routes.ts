import { index, layout, route, type RouteConfig } from "@react-router/dev/routes";

export default [
  index("./routes/login.tsx"),
  layout("./layout.tsx", [
    route("deployments", "./routes/deployments.tsx"),
    route("new", "./routes/new.tsx"),
    route("_deployments_test", "./routes/_deployments_test.tsx"),
  ]),
] satisfies RouteConfig;
