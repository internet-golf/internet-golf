import { Link, Outlet, useMatch } from "react-router";
import { Drawer, Flex, Space, theme, Typography } from "antd";
import {
  DeploymentUnitOutlined,
  AppstoreAddOutlined,
  BarsOutlined,
  MenuOutlined,
} from "@ant-design/icons";
import ColorScheme from "./components/ColorScheme";
import { useRef, useState, type ReactNode } from "react";

function HeaderLink({ children, to, icon }: { children: ReactNode; to: string; icon?: ReactNode }) {
  const isCurrent = useMatch(to);
  return (
    <Space align="baseline" size="small">
      {icon}
      <Link to={to}>
        <Typography.Title
          level={5}
          style={{
            margin: 0,
            textDecoration: isCurrent ? "none" : "underline",
          }}
        >
          {children}
        </Typography.Title>
      </Link>
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
    <Flex
      align={vertical ? "start" : "center"}
      vertical={vertical}
      gap={vertical ? "small" : "large"}
    >
      <HeaderLink to="/deployments" icon={<DeploymentUnitOutlined />}>
        Deployments
      </HeaderLink>
      <HeaderLink to="/domains" icon={<BarsOutlined />}>
        Domains
      </HeaderLink>
      <HeaderLink to="/permissions" icon={<BarsOutlined />}>
        Permissions
      </HeaderLink>
      <div className={vertical ? "" : "ml-auto"}>
        <HeaderLink to="/new" icon={<AppstoreAddOutlined />}>
          Deploy New Site
        </HeaderLink>
      </div>
      {vertical && <StatusLink />}
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
        <div className="flex flex-col [&>div]:py-2 [&>div]:px-4 [&>div]:mx-auto [&>div]:w-5xl [&>div]:max-w-full [&>div]:md:max-w-[95%]">
          <div className="md:mb-1 mb-3 mt-3">
            <Flex align="baseline" gap="middle">
              <button className="contents md:hidden" onClick={() => setDrawerOpen((o) => !o)}>
                <MenuOutlined style={{ fontSize: token.fontSizeHeading4 }} />
              </button>
              <Link to="/">
                <Typography.Title level={2} style={{ margin: 0 }}>
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
          styles={{ root: { position: "absolute" }, body: { padding: "10px" } }}
          placement="left"
          size={275}
          closable={false}
          onClose={() => setDrawerOpen(false)}
          open={drawerOpen}
          getContainer={drawerContainer.current ?? false}
        >
          <HeaderLinks vertical />
        </Drawer>
      </div>
      <ColorScheme.Light>
        <div className="w-5xl max-w-[95%] mx-auto py-2">
          <Outlet />
        </div>
      </ColorScheme.Light>
    </>
  );
}
