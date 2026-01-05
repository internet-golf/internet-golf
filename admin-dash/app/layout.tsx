import { Outlet } from "react-router";
import { Divider, Drawer, Flex, Space, theme, Typography } from "antd";
import {
  DeploymentUnitOutlined,
  AppstoreAddOutlined,
  BarsOutlined,
  MenuOutlined,
} from "@ant-design/icons";
import ColorScheme from "./components/ColorScheme";
import { useRef, useState, type ReactNode } from "react";

function HeaderLink({ children, to, icon }: { children: ReactNode; to: string; icon?: ReactNode }) {
  return (
    <Space align="baseline" size="middle">
      {icon}
      <a href={to}>
        <Typography.Title level={5} style={{ margin: 0, textDecoration: "underline" }}>
          {children}
        </Typography.Title>
      </a>
    </Space>
  );
}

function HeaderLinks({ vertical }: { vertical?: boolean }) {
  return (
    <Space
      align={vertical ? "start" : "center"}
      vertical={vertical}
      size={vertical ? "small" : "large"}
    >
      <HeaderLink to="/deployments" icon={<DeploymentUnitOutlined />}>
        Deployments
      </HeaderLink>
      <HeaderLink to="/domains" icon={<BarsOutlined />}>
        Domains
      </HeaderLink>
      <HeaderLink to="/new" icon={<AppstoreAddOutlined />}>
        Deploy New Site
      </HeaderLink>
    </Space>
  );
}

export default function LayoutComponent() {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const drawerContainer = useRef<HTMLDivElement | null>(null);
  const { token } = theme.useToken();
  return (
    <>
      <ColorScheme.Dark>
        <div className="flex flex-col [&>div]:py-2 [&>div]:px-4 [&>div]:mx-auto [&>div]:w-5xl [&>div]:max-w-full">
          <div className="mb-2 mt-3">
            <Flex align="baseline" gap="middle">
              <button className="contents md:hidden" onClick={() => setDrawerOpen((o) => !o)}>
                <MenuOutlined style={{ fontSize: token.fontSizeHeading4 }} />
              </button>
              <Typography.Title level={2} style={{ margin: 0 }}>
                Internet Golf
              </Typography.Title>
              <div className="ml-auto bg-actual-gray rounded-lg self-center p-2 lg:-mr-4">
                <HeaderLink
                  to="/status"
                  icon={<div className="w-2 h-2 mb-0.5 bg-green-600 rounded-full"> </div>}
                >
                  Server Status
                </HeaderLink>
              </div>
            </Flex>
          </div>
          <div className="hidden md:block bg-actual-gray rounded-t-lg">
            <HeaderLinks />
          </div>
        </div>
      </ColorScheme.Dark>
      <div ref={drawerContainer} className="absolute left-0 h-full w-full">
        <Drawer
          styles={{ root: { position: "absolute" }, body: { padding: "10px" } }}
          placement="left"
          size={350}
          closable={false}
          onClose={() => setDrawerOpen(false)}
          open={drawerOpen}
          getContainer={drawerContainer.current ?? false}
        >
          <HeaderLinks vertical />
        </Drawer>
      </div>
      <ColorScheme.Light>
        <div className="w-5xl max-w-full mx-auto py-2">
          <Outlet />
        </div>
      </ColorScheme.Light>
    </>
  );
}
