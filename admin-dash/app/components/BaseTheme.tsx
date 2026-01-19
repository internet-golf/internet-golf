import { ConfigProvider } from "antd";
import type { ReactNode } from "react";

const roundedBorderRadiusValues = {
  borderRadiusSM: 4,
  borderRadius: 6,
};

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
          paddingMD: 9,
          padding: 10,
          paddingLG: 16,

          // make borders a squarer and darker
          colorBorder: "#bbb",
          colorBorderSecondary: "#bbb",
          borderRadius: 0,
          lineWidth: 1,
        },
        components: {
          // grant buttons and inputs an exemption from square borders
          Button: roundedBorderRadiusValues,
          Input: roundedBorderRadiusValues,
          InputNumber: roundedBorderRadiusValues,
          Segmented: roundedBorderRadiusValues,
          Select: roundedBorderRadiusValues,
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
