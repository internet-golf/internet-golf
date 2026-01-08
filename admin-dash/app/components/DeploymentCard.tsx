import {
  EditOutlined,
  ExportOutlined,
  FolderOpenFilled,
  GithubOutlined,
  ReloadOutlined,
} from "@ant-design/icons";
import { Button, Card, Flex, Typography } from "antd";
import type { ReactNode } from "react";
import type { GetDeploymentResponse } from "~/api-calls/generated/golfComponents";

function ExternalSourceLink({
  externalSource,
  externalSourceType,
}: Pick<GetDeploymentResponse, "externalSource" | "externalSourceType">) {
  if (externalSourceType === "GithubRepo") {
    return (
      <Flex gap="small" align="center">
        <GithubOutlined />
        <a target="_blank" href={`https://github.com/${externalSource}`}>
          {externalSource}
        </a>
      </Flex>
    );
  }
  return null;
}

const TypeLabel = ({ icon, children }: { icon: ReactNode; children: ReactNode }) => (
  <Flex align="center" gap="small">
    {icon}
    {children}
  </Flex>
);

export function DeploymentCard({ deployment }: { deployment: GetDeploymentResponse }) {
  const type = {
    Alias: <TypeLabel icon={<ReloadOutlined />}>Alias</TypeLabel>,
    Empty: (
      <TypeLabel
        icon={<span className="rounded-full border-2 border-dashed w-3 h-3 inline-block"> </span>}
      >
        No Content Yet
      </TypeLabel>
    ),
    StaticSite: <TypeLabel icon={<FolderOpenFilled />}>Static Site</TypeLabel>,
  }[deployment.type];

  const actionButtons = [<Button icon={<EditOutlined />}>Edit</Button>];

  if (deployment.type !== "Empty") {
    actionButtons.push(
      <a href={`//${deployment.url}`} target="_blank">
        <Button type="link" icon={<ExportOutlined />} iconPlacement="end">
          Visit
        </Button>
      </a>,
    );
  }

  return (
    <Card
      title={
        <Flex vertical>
          <Typography.Text type="secondary">{type}</Typography.Text>
          <Typography.Title underline level={5} style={{ margin: 0 }}>
            {deployment.url}
          </Typography.Title>
        </Flex>
      }
      styles={{
        root: { width: "100%", display: "flex", flexDirection: "column" },
        title: { paddingTop: 12, paddingBottom: 12 },
        body: { paddingTop: 12, paddingBottom: 12, minHeight: 0, flexGrow: "1" },
      }}
      actions={actionButtons}
    >
      <div className="flex flex-col justify-evenly pt-1 gap-1 h-full">
        {deployment.type === "Alias" && (
          <Typography.Text>
            {deployment.redirect ? "Redirect" : "Alias"} to <strong>{deployment.aliasedTo}</strong>
          </Typography.Text>
        )}
        {!!deployment.externalSource && <ExternalSourceLink {...deployment} />}
        {!!deployment.tags?.length && (
          <Typography.Text type="secondary">
            Tags: <Typography.Text>{(deployment.tags ?? []).join(", ")}</Typography.Text>
          </Typography.Text>
        )}
        <Typography.Text type="secondary" title={new Date(deployment.updatedAt).toLocaleString()}>
          Last updated{" "}
          {new Date(deployment.updatedAt).toLocaleDateString("en-US", { dateStyle: "long" })}
        </Typography.Text>
      </div>
    </Card>
  );
}
