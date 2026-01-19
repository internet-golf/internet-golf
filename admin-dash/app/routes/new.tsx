import { useState } from "react";
import { Form, Radio, Typography } from "antd";
import type { Route } from "./+types/new";
import { StaticSiteForm } from "~/components/forms/StaticSiteForm";
import { AliasForm } from "~/components/forms/AliasForm";
import { EmptyDeploymentForm } from "~/components/forms/EmptyDeploymentForm";
import type { GetDeploymentResponse } from "~/api-calls/generated/golfComponents";

export function meta({}: Route.MetaArgs) {
  return [{ title: "Deploy New Site" }];
}

type DeploymentType = GetDeploymentResponse["type"];

export default function NewDeployment() {
  const [deploymentType, setDeploymentType] = useState<DeploymentType | null>(null);

  return (
    <>
      <Typography.Title level={2}>Deploy New Site</Typography.Title>

      <div className="flex flex-col md:flex-row gap-8">
        <div className="md:w-1/4 shrink-0">
          <div className="border border-gray-200 p-4">
            <Typography.Title level={5} style={{ margin: 0 }}>
              Deployment Type
            </Typography.Title>
            <div className="mt-3">
              <Radio.Group
                vertical
                value={deploymentType}
                onChange={(e) => setDeploymentType(e.target.value)}
                options={[
                  { value: "StaticSite", label: "Static Site" },
                  { value: "Alias", label: "Alias" },
                  { value: "Empty", label: "Empty (No content yet)" },
                ]}
              />
            </div>
          </div>
        </div>

        <div className="w-full flex-1 max-w-2xl">
          {deploymentType === null && (
            <div className="border border-gray-200 p-4 text-center text-gray-500">
              <Typography.Text type="secondary">Select a deployment type.</Typography.Text>
            </div>
          )}
          {deploymentType === "StaticSite" && <StaticSiteForm />}
          {deploymentType === "Alias" && <AliasForm />}
          {deploymentType === "Empty" && <EmptyDeploymentForm />}
        </div>
      </div>
    </>
  );
}
