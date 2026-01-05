import type { Route } from "./+types/login";

import { Button, Flex, Form, Input, Modal, Typography } from "antd";

import { useState } from "react";

import ColorScheme from "~/components/ColorScheme";

export function meta({}: Route.MetaArgs) {
  return [{ title: "Login" }];
}

export default function Login() {
  const [isAdvancedModalOpen, setIsAdvancedModalOpen] = useState(false);
  const [isHelpModalOpen, setIsHelpModalOpen] = useState(false);
  return (
    <ColorScheme.Dark>
      <div className="w-screen h-screen flex flex-col justify-center items-center relative">
        {/* font size override just for the initial title - `style` overrides are otherwise discouraged */}
        <Typography.Title style={{ fontSize: "3em" }} className="text-center">
          Internet Golf
          <br />
          Dashboard
        </Typography.Title>
        <Flex gap="small">
          <Input.Password placeholder="Login Token" />
          <Button>Login</Button>
        </Flex>
        <div className="fixed right-4 md:right-8 bottom-4 md:bottom-8 flex gap-2">
          <Button onClick={() => setIsAdvancedModalOpen(true)}>Advanced...</Button>
          <Button onClick={() => setIsHelpModalOpen(true)}>?</Button>
        </div>
        <Modal
          title="Advanced Settings"
          open={isAdvancedModalOpen}
          okText="Save"
          onOk={() => setIsAdvancedModalOpen(false)}
          onCancel={() => setIsAdvancedModalOpen(false)}
        >
          <Form layout="vertical">
            <Form.Item label="Manually Set API URL">
              <Input placeholder="http://localhost:8888/" />
            </Form.Item>
          </Form>
        </Modal>
        <Modal
          title="?"
          open={isHelpModalOpen}
          okText="Close"
          onCancel={() => setIsHelpModalOpen(false)}
          onOk={() => setIsHelpModalOpen(false)}
          // this skips the "cancel" button in the footer
          footer={(_, { OkBtn }) => <OkBtn />}
        >
          <Typography.Paragraph>Welcome to the Internet Golf Dashboard!</Typography.Paragraph>
          <Typography.Paragraph>
            You can obtain a login token for your server by running the{" "}
            <Typography.Text code>golf create-token</Typography.Text> command. See the{" "}
            <a href="https://github.com/internet-golf/internet-golf">repo</a> for more details.
          </Typography.Paragraph>
        </Modal>
      </div>
    </ColorScheme.Dark>
  );
}
