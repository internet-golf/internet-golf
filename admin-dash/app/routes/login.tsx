import type { Route } from "./+types/login";

import { Button, Form, Input, Modal, Typography } from "antd";

import { useState } from "react";
import { redirect, useNavigate } from "react-router";
import { getGolfAuthToken, setGolfAuthToken } from "~/api-calls/session";

import ColorScheme from "~/components/ColorScheme";

export function meta({}: Route.MetaArgs) {
  return [{ title: "Login" }];
}

export function clientLoader() {
  // if the user already has a token stored, send them onward. if the token is
  // invalid, then golfFetcher will receive a 401 Unauthorized response and send
  // them right back here; however, by attempting to forward them to the next
  // page by default prioritizes the common/successful case
  if (getGolfAuthToken()) {
    return redirect("/deployments");
  }
}

export default function Login() {
  let navigate = useNavigate();

  const [isHelpModalOpen, setIsHelpModalOpen] = useState(false);

  const login = async ({ token }: { token: string }) => {
    if (!token) {
      console.error("no token");
      return;
    }
    console.log(token);
    setGolfAuthToken(token);
    navigate("/deployments");
  };

  return (
    <ColorScheme.Dark>
      <div className="w-screen h-screen flex flex-col justify-center items-center relative px-4">
        <Typography.Title style={{ fontSize: "3em", textAlign: "center", textWrap: "balance" }}>
          Internet Golf Admin
        </Typography.Title>
        <Form onFinish={login}>
          <div className="flex justify-center items-center gap-4 max-w-full">
            <Form.Item
              name="token"
              rules={[
                {
                  required: true,
                  message: "Please input an authentication token.",
                },
              ]}
            >
              <Input.Password
                name="token"
                style={{ width: "350px" }}
                placeholder="Enter authentication token..."
              />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit">
                Login
              </Button>
            </Form.Item>
          </div>
        </Form>
        <div className="fixed right-4 md:right-6 bottom-4 md:bottom-6 flex gap-2">
          <Button onClick={() => setIsHelpModalOpen(true)}>?</Button>
        </div>
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
