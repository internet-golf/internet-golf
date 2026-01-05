import { ConfigProvider } from "antd";
import type { ReactNode } from "react";

/**
 * Style defaults, applied in root.tsx.
 */
export default function BaseTheme({ children }: { children: ReactNode }) {
  return (
    <ConfigProvider
      // the modal body looks cramped to me without some extra spacing around it
      modal={{ styles: { body: { margin: "12px 0" } } }}
      theme={{
        token: {
          // this makes the <Space> component (and the <Flex> component?)
          // give less space between items by default
          paddingXS: 6,
          paddingSM: 7,
          padding: 8,
        },
        components: {
          Modal: {
            // i am not that enthusiastic about the swoosh animation for this thing
            motionDurationFast: "0",
            motionDurationMid: "0",
            motionDurationSlow: "0",
          },
        },
      }}
    >
      {children}
    </ConfigProvider>
  );
}
