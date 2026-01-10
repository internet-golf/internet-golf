import { ConfigProvider, theme } from "antd";
import type { ReactNode } from "react";

/**
 * Provides color scheme-specific CSS config, and also a container for the
 * `children` with the logical background.
 */
function ColorSchemeApplicator({
  children,
  mode,
}: {
  children: ReactNode;
  mode: "dark" | "light";
}) {
  return (
    // note that this inherits some global config from the BaseTheme (since
    // that's used in root.tsx)
    <ConfigProvider
      theme={{
        algorithm: mode === "dark" ? theme.darkAlgorithm : theme.defaultAlgorithm,
        // imo, the dark mode text looks better when it's brighter than the antd default
        token: { colorText: mode === "dark" ? "#fdfdfd" : "black" },
      }}
      // this class is used to set the icon color in app.css
      iconPrefixCls={mode === "dark" ? "dark-mode-icon" : "light-mode-icon"}
      divider={{
        // for some reason the divider is hard to see in dark mode by default
        style: { borderColor: mode === "dark" ? "var(--ant-color-text)" : undefined },
      }}
    >
      <div className={mode === "dark" ? "bg-black" : "bg-white"}>{children}</div>
    </ConfigProvider>
  );
}

/**
 * Wrappers that apply color schemes.
 *
 * @example
 * <ColorScheme.Dark>
 *   <Typography.Paragraph>Hi!</Typography.Paragraph>
 * </ColorScheme.Dark>
 */
const ColorScheme = {
  Dark: ({ children }: { children: ReactNode }) => (
    <ColorSchemeApplicator mode="dark">{children}</ColorSchemeApplicator>
  ),
  Light: ({ children }: { children: ReactNode }) => (
    <ColorSchemeApplicator mode="light">{children}</ColorSchemeApplicator>
  ),
};

export default ColorScheme;
