import { Outlet } from "react-router";
import { Divider, Flex, Space, Typography } from "antd";
import { DeploymentUnitOutlined, AppstoreAddOutlined, BarsOutlined } from "@ant-design/icons";
import ColorScheme from "./components/ColorScheme";
import type { ReactNode } from "react";

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

export default function LayoutComponent() {
  return (
    <>
      <ColorScheme.Dark>
        <div className="p-4">
          <Flex vertical>
            {/* overriding the default margins for this title only  */}
            <Typography.Title style={{ margin: "0 0 12px" }}>Internet Golf</Typography.Title>
            <Space separator={<Divider vertical />} align="center">
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
          </Flex>
        </div>
      </ColorScheme.Dark>
      <ColorScheme.Light>
        <Outlet />
      </ColorScheme.Light>
    </>
  );
}
