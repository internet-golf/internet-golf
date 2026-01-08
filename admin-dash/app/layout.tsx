import { Link, Outlet, useMatch } from "react-router";
import { Drawer, Flex, Space, theme, Typography } from "antd";
import {
  DeploymentUnitOutlined,
  AppstoreAddOutlined,
  BarsOutlined,
  MenuOutlined,
  ClusterOutlined,
  CloudUploadOutlined,
  CloudServerOutlined,
} from "@ant-design/icons";
import ColorScheme from "./components/ColorScheme";
import { useRef, useState, type ReactNode } from "react";

function HeaderLink({
  children,
  to,
  icon,
  iconPlacement = "left",
}: {
  children: ReactNode;
  to: string;
  icon?: ReactNode;
  iconPlacement?: "left" | "right";
}) {
  const isCurrent = useMatch(to);
  return (
    <Space align="baseline" size="small">
      {iconPlacement === "left" && <div className="w-4">{icon}</div>}
      <Link to={to}>
        <Typography.Title
          level={5}
          underline={!isCurrent}
          style={{
            margin: 0,
            cursor: isCurrent ? "default" : "pointer",
          }}
        >
          {children}
        </Typography.Title>
      </Link>
      {iconPlacement === "right" && <div className="w-4">{icon}</div>}
    </Space>
  );
}

function StatusLink() {
  return (
    <HeaderLink
      to="/status"
      icon={<div className="w-2 h-2 mx-[3px] mb-px bg-(--ant-color-success) rounded-full"> </div>}
    >
      Server Status
    </HeaderLink>
  );
}

function HeaderLinks({ vertical }: { vertical?: boolean }) {
  return (
    <Flex align={vertical ? "start" : "center"} vertical={vertical} gap="small">
      <HeaderLink to="/deployments" icon={<CloudServerOutlined />}>
        Deployments
      </HeaderLink>
      <HeaderLink to="/domains" icon={null}>
        Domains
      </HeaderLink>
      <HeaderLink to="/permissions" icon={null}>
        Permissions
      </HeaderLink>
      <div className={vertical ? "" : "ml-auto"}>
        <HeaderLink
          to="/new"
          iconPlacement={vertical ? "left" : "right"}
          icon={<CloudUploadOutlined />}
        >
          Deploy New Site
        </HeaderLink>
      </div>
      {vertical ? <StatusLink /> : null}
    </Flex>
  );
}

export default function LayoutComponent() {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const drawerContainer = useRef<HTMLDivElement | null>(null);
  const { token } = theme.useToken();

  return (
    <>
      <ColorScheme.Dark>
        <div className="flex flex-col [&>div]:py-2 [&>div]:px-4 [&>div]:mx-auto [&>div]:w-6xl [&>div]:max-w-full [&>div]:md:max-w-[95%]">
          <div className="mb-3 mt-4 md:my-4">
            <Flex align="center" gap="middle">
              <button
                className="block md:hidden mt-2 mr-1"
                onClick={() => setDrawerOpen((o) => !o)}
              >
                <MenuOutlined style={{ fontSize: token.fontSizeHeading3 }} />
              </button>
              <Link to="/">
                <Typography.Title level={1} style={{ margin: 0 }}>
                  Internet Golf
                </Typography.Title>
              </Link>
              <div className="hidden md:block ml-auto self-center p-2">
                <StatusLink />
              </div>
            </Flex>
          </div>
          <div className="hidden md:block bg-actual-gray rounded-lg mb-3">
            <HeaderLinks />
          </div>
        </div>
      </ColorScheme.Dark>
      <div ref={drawerContainer} className="absolute left-0 w-full h-full">
        <Drawer
          styles={{ root: { position: "absolute" }, body: { padding: "12px 16px" } }}
          placement="left"
          size={275}
          closable={false}
          onClose={() => setDrawerOpen(false)}
          open={drawerOpen}
          getContainer={drawerContainer.current ?? false}
        >
          <ColorScheme.Light>
            <HeaderLinks vertical />
          </ColorScheme.Light>
        </Drawer>
      </div>
      <ColorScheme.Light>
        <div className="w-6xl max-w-[95%] mx-auto py-4">
          <Outlet />
        </div>
      </ColorScheme.Light>
    </>
  );
}
