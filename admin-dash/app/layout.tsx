import { Link, Outlet, useMatch } from "react-router";
import { Drawer, Flex, theme, Typography } from "antd";
import { MenuOutlined, CloudUploadOutlined } from "@ant-design/icons";
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
    <Flex align="center" gap="small">
      {!!icon && iconPlacement === "left" && (
        <div className="w-4 h-4 mt-1 flex justify-center items-center">{icon}</div>
      )}
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
      {!!icon && iconPlacement === "right" && (
        <div className="w-4 h-4 mt-1 flex justify-center items-center">{icon}</div>
      )}
    </Flex>
  );
}

function StatusLink({ iconPlacement = "right" }: { iconPlacement?: "left" | "right" }) {
  return null;
  return (
    <HeaderLink
      to="/status"
      icon={<div className="w-2 h-2 mx-[3px] mb-px bg-green-500 rounded-full"> </div>}
      iconPlacement={iconPlacement}
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
      <HeaderLink to="/deployments" icon={null}>
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
      {vertical ? <StatusLink iconPlacement="left" /> : null}
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
        <div className="flex flex-col">
          <div className="px-4 md:px-0 mb-3 mt-4 md:mt-5 md:mb-3 py-2 mx-auto w-6xl max-w-full md:max-w-[95%]">
            <Flex align="center" gap="middle">
              <button
                className="block md:hidden mt-2 mr-1"
                onClick={() => setDrawerOpen((o) => !o)}
              >
                <MenuOutlined style={{ fontSize: token.fontSizeHeading4 }} />
              </button>
              <Link to="/">
                <div className="hidden md:block">
                  <Typography.Title level={1} style={{ margin: 0 }}>
                    Internet Golf Dashboard
                  </Typography.Title>
                </div>
                <div className="md:hidden block">
                  <Typography.Title level={2} style={{ margin: 0 }}>
                    Internet Golf Dashboard
                  </Typography.Title>
                </div>
              </Link>
              <div className="hidden md:block ml-auto self-center p-2">
                <StatusLink />
              </div>
            </Flex>
          </div>
          <div className="hidden md:block py-2 mb-1 bg-[#1c1c1c]">
            <div className="mx-auto w-6xl max-w-full md:max-w-[95%]">
              <HeaderLinks />
            </div>
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
        <div className="w-6xl max-w-[95%] mx-auto py-6">
          <Outlet />
        </div>
      </ColorScheme.Light>
    </>
  );
}
