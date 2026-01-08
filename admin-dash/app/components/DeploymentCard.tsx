import {
  EditOutlined,
  ExportOutlined,
  FolderOpenFilled,
  GithubOutlined,
  ReloadOutlined,
} from "@ant-design/icons";
import { Button, Card, Flex, Tag, theme, Typography } from "antd";
import type { ReactNode } from "react";
import type { GetDeploymentResponse } from "~/api-calls/generated/golfComponents";
import { allowBreakingOnDots } from "~/utils/utils";

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
  const { token } = theme.useToken();

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
      styles={{
        root: { width: "100%", display: "flex", flexDirection: "column" },
        title: { padding: 0 },
        header: { padding: "12px 16px" },
        body: { padding: 16, minHeight: 0, flexGrow: "1" },
      }}
      actions={actionButtons}
    >
      <Flex gap="middle">
        <div className="h-24 w-32 rounded-lg bg-gray-100">
          <img src={deployment.meta.image || "/full-web-icon.svg"} />
        </div>
        <Flex gap="small" vertical>
          {deployment.type === "Alias" ? (
            <Typography.Title level={5}>Alias to {deployment.aliasedTo}</Typography.Title>
          ) : (
            <>
              {deployment.meta.title ? (
                <Typography.Text style={{ fontSize: token.fontSizeHeading5 }}>
                  <strong>{deployment.meta.title}</strong> ({allowBreakingOnDots(deployment.url)})
                </Typography.Text>
              ) : (
                <Typography.Text strong style={{ fontSize: token.fontSizeHeading5 }}>
                  {allowBreakingOnDots(deployment.url)}
                </Typography.Text>
              )}
              {!!deployment.meta.description && (
                <Typography.Text>{deployment.meta.description}</Typography.Text>
              )}
              {!!deployment.externalSource && <ExternalSourceLink {...deployment} />}
            </>
          )}
          {/* <Typography.Text type="secondary">
            Last updated{" "}
            {new Date(deployment.updatedAt).toLocaleDateString("en-US", { dateStyle: "long" })}
          </Typography.Text> */}
          <Tag style={{ alignSelf: "flex-start", marginTop: "4px" }}>{type}</Tag>
        </Flex>
      </Flex>
    </Card>
  );
}
