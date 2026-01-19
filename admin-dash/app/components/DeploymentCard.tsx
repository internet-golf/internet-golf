import {
  DeleteOutlined,
  FolderOpenFilled,
  GithubOutlined,
  LinkOutlined,
  ReloadOutlined,
} from "@ant-design/icons";
import { Button, Card, Flex, Popconfirm, Tag, theme, Typography } from "antd";
import { useState } from "react";
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
  if (externalSourceType === "Github") {
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

export function DeploymentCard({
  deployment,
  onDelete,
}: {
  deployment: GetDeploymentResponse;
  onDelete: () => Promise<void>;
}) {
  const { token } = theme.useToken();
  const [isDeleting, setIsDeleting] = useState(false);

  const name = deployment.name || deployment.meta.title;

  const handleDelete = async () => {
    setIsDeleting(true);
    await onDelete();
    setIsDeleting(false);
  };

  return (
    <Card
      styles={{
        root: { width: "100%", display: "flex", flexDirection: "column" },
        body: { padding: 16, minHeight: 0, flexGrow: "1" },
      }}
    >
      <Popconfirm
        title="Delete deployment"
        description={`Are you sure you want to delete ${deployment.url}?`}
        onConfirm={handleDelete}
        okText="Delete"
        okButtonProps={{ danger: true }}
        cancelText="Cancel"
      >
        <Button
          icon={<DeleteOutlined />}
          variant="outlined"
          loading={isDeleting}
          style={{ float: "right", marginLeft: 12 }}
        />
      </Popconfirm>
      <Flex gap="large">
        <div className="h-24 w-32 bg-gray-100">
          <img
            className="h-full w-auto mx-auto"
            src={deployment.meta.image || "/internet-icon.svg"}
          />
        </div>
        <Flex gap="small" vertical justify="flex-start">
          {deployment.type === "Alias" ? (
            // using token.fontSizeHeading5 to get a standard font size instead
            // of using Typography.Title, to avoid getting the standard margins
            // for an h5
            <Typography.Text strong style={{ fontSize: token.fontSizeHeading5 }}>
              {deployment.redirect ? "Redirect" : "Alias"} to {deployment.aliasedTo}
            </Typography.Text>
          ) : (
            <Typography.Text strong style={{ fontSize: token.fontSizeHeading5 }}>
              {name || allowBreakingOnDots(deployment.url)}
            </Typography.Text>
          )}
          <Typography.Text>
            <a href={`//${deployment.url}`} target="_blank" rel="nofollow noreferrer">
              <LinkOutlined /> {allowBreakingOnDots(deployment.url)}
            </a>
          </Typography.Text>
          {/* TODO: i do not know why this marginTop is needed to make things look spaced evenly */}
          <Flex gap="small" align="center" style={{ marginTop: 6 }}>
            <Tag style={{ alignSelf: "flex-start" }}>
              <TypeLabel {...deployment} />
            </Tag>
            {deployment.externalSource && <ExternalSourceTag {...deployment} />}
          </Flex>
        </Flex>
      </Flex>
    </Card>
  );
}
