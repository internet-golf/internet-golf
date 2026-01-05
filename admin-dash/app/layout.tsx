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
    <Space align="baseline" size="small">
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
      separator={vertical ? null : <Divider vertical />}
      align={vertical ? "start" : "center"}
      vertical={vertical}
    >
      <HeaderLink to="/deployments" icon={<DeploymentUnitOutlined />}>
        Deployments
      </HeaderLink>
      <HeaderLink to="/domains" icon={<BarsOutlined />}>
        Domains
      </HeaderLink>
      <HeaderLink
        to="/status"
        icon={<div className="w-2 h-2 mb-0.5 bg-green-600 rounded-full"> </div>}
      >
        Server Status
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
        <div className="p-4">
          <Flex vertical>
            <div className="m-0 md:mb-2">
              <Flex align="baseline" gap="middle">
                <button className="contents md:hidden" onClick={() => setDrawerOpen((o) => !o)}>
                  <MenuOutlined style={{ fontSize: token.fontSizeHeading4 }} />
                </button>
                <Typography.Title level={2} style={{ margin: 0 }}>
                  Internet Golf
                </Typography.Title>
              </Flex>
            </div>
            <div className="hidden md:contents">
              <HeaderLinks />
            </div>
          </Flex>
        </div>
      </ColorScheme.Dark>
      <div ref={drawerContainer} className="absolute left-0 h-full w-full">
        <Drawer
          styles={{ root: { position: "absolute" }, body: { padding: "16px" } }}
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
        <Outlet />
      </ColorScheme.Light>
    </>
  );
}
