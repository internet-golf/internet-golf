import {
  EditOutlined,
  ExportOutlined,
  FolderOpenFilled,
  GithubOutlined,
  LinkOutlined,
  ReloadOutlined,
} from "@ant-design/icons";
import { Button, Card, Flex, Tag, theme, Typography } from "antd";
import type { ReactNode } from "react";
import type { GetDeploymentResponse } from "~/api-calls/generated/golfComponents";
import { allowBreakingOnDots } from "~/utils/utils";

const TypeLabel = ({
  type,
  redirect,
}: Pick<GetDeploymentResponse, "type"> & { redirect?: boolean }) => {
  const typeConfig = {
    Alias: { icon: <ReloadOutlined />, label: redirect ? "Redirect" : "Alias" },
    Empty: {
      icon: <span className="rounded-full border-2 border-dashed w-3 h-3 inline-block"> </span>,
      label: "No Content Yet",
    },
    StaticSite: { icon: <FolderOpenFilled />, label: "Static Site" },
  }[type];

  if (!typeConfig) return null;

  return (
    <Flex align="center" gap="small">
      {typeConfig.icon}
      {typeConfig.label}
    </Flex>
  );
};

function ExternalSourceTag({
  externalSource,
  externalSourceType,
}: Pick<GetDeploymentResponse, "externalSource" | "externalSourceType">) {
  if (externalSourceType === "GithubRepo") {
    return (
      <Tag>
        <a target="_blank" href={`https://github.com/${externalSource}`}>
          <GithubOutlined /> <span style={{ textDecoration: "underline" }}>Github</span>
        </a>
      </Tag>
    );
  }
  return null;
}

export function DeploymentCard({ deployment }: { deployment: GetDeploymentResponse }) {
  const { token } = theme.useToken();

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

  const name = deployment.name || deployment.meta.title;

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
      <Flex gap="large">
        <div className="h-24 w-32 rounded-lg overflow-clip bg-gray-100">
          <img src={deployment.meta.image || "/full-web-icon.svg"} />
        </div>
        <Flex gap="middle" vertical justify="flex-start">
          {deployment.type === "Alias" ? (
            <Typography.Title level={5}>
              {deployment.redirect ? "Redirect" : "Alias"} to {deployment.aliasedTo}
            </Typography.Title>
          ) : (
            // using token.fontSizeHeading5 to get a standard font size instead
            // of using Typography.Title, to avoid getting the standard margins
            // for an h5
            <>
              {name ? (
                <>
                  <Typography.Text strong style={{ fontSize: token.fontSizeHeading5 }}>
                    {name}
                  </Typography.Text>
                  {/* <Typography.Text>
                    <a href={`//${deployment.url}`} target="_blank" rel="nofollow noreferrer">
                      <LinkOutlined /> {allowBreakingOnDots(deployment.url)}
                    </a>
                  </Typography.Text> */}
                </>
              ) : (
                <Typography.Text strong style={{ fontSize: token.fontSizeHeading5 }}>
                  {allowBreakingOnDots(deployment.url)}
                </Typography.Text>
              )}
            </>
          )}
          {!!deployment.meta.description && (
            <Typography.Text>{deployment.meta.description}</Typography.Text>
          )}
          {/* TODO: i do not know why this marginTop is needed to make things look spaced evenly */}
          <Flex gap="small" align="center" style={{ marginTop: 6 }}>
            {deployment.externalSource && <ExternalSourceTag {...deployment} />}
            <Tag style={{ alignSelf: "flex-start" }}>
              <TypeLabel {...deployment} />
            </Tag>
          </Flex>
        </Flex>
      </Flex>
    </Card>
  );
}
