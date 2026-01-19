import type { Route } from "./+types/login";

import { Button, Form, Input, Modal, Typography } from "antd";

import { useState } from "react";
import { redirect, useNavigate } from "react-router";
import { getGolfAuthToken, setGolfAuthToken } from "~/api-calls/session";

import ColorScheme from "~/components/ColorScheme";

interface LoginFormValues {
  token: string;
}

const FormItem = Form.Item<LoginFormValues>;

export function meta({}: Route.MetaArgs) {
  return [{ title: "Internet Golf Admin Dashboard", description: "Manage your server." }];
}

export function clientLoader() {
  // if the user already has a token stored, send them onward. if the token has
  // been invalidated since it was entered, then golfFetcher will receive a 401
  // Unauthorized response the next time it tries to make a request and then
  // will send them right back here, so this might not "work"; however,
  // attempting to forward them to the next page by default prioritizes the
  // common, successful case
  if (getGolfAuthToken()) {
    return redirect("/deployments");
  }
}

export default function Login() {
  let navigate = useNavigate();

  const [isHelpModalOpen, setIsHelpModalOpen] = useState(false);

  const login = async ({ token }: LoginFormValues) => {
    if (!token) {
      console.error("no token");
      return;
    }
    // TODO: some kind of call to like a "/validate-auth" endpoint should be
    // made here before accepting and storing the token and forwarding the user;
    // it would speed up the process of finding out the token is wrong
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
          <div className="flex flex-col gap-1 justify-center items-center max-w-full">
            <FormItem
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
            </FormItem>
            <FormItem>
              <Button type="primary" htmlType="submit">
                Login
              </Button>
            </FormItem>
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
