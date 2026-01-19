/**
 * This contains raw materials that the forms for creating specific deployments
 * should use. They should extend BaseDeploymentValues to create their form
 * values type; put CommonDeploymentFields and AdvancedFields in the bodies of
 * their forms; and call createDeployment from this file to use the values from
 * those fields.
 */

import { Checkbox, Col, Collapse, Form, Input, Row, Select, Space, Tooltip } from "antd";
import { InfoCircleOutlined } from "@ant-design/icons";
import { createDeployment as createDeploymentApiCall } from "~/api-calls/generated/golfComponents";

export interface BaseDeploymentValues {
  url: string;
  name?: string;
  externalSource?: string;
  externalSourceType?: string;
  preserveExternalPath?: boolean;
}

const FormItem = Form.Item<BaseDeploymentValues>;

export function CommonDeploymentFields() {
  return (
    <>
      <FormItem label="URL" name="url" rules={[{ required: true, message: "Please enter a URL" }]}>
        <Input placeholder="mysite.mydomain.com" />
      </FormItem>

      <Row gutter={16}>
        <Col xs={24} md={12}>
          <FormItem label="Name (optional)" name="name">
            <Input placeholder="My Site" />
          </FormItem>
        </Col>
        <Col xs={24} md={12}>
          <FormItem label="Repo (optional)">
            <Space.Compact block>
              <FormItem name="externalSourceType">
                <Select
                  style={{ width: 100 }}
                  placeholder="Source"
                  allowClear
                  options={[{ label: "Github", value: "Github" }]}
                />
              </FormItem>
              <FormItem name="externalSource" style={{ width: "100%" }}>
                <Input placeholder="username/reponame" />
              </FormItem>
            </Space.Compact>
          </FormItem>
        </Col>
      </Row>
    </>
  );
}

export function AdvancedFields() {
  return (
    <Collapse
      ghost
      items={[
        {
          key: "advanced",
          label: "Advanced",
          children: (
            <FormItem name="preserveExternalPath" valuePropName="checked">
              <Checkbox>
                Preserve internal path{" "}
                <Tooltip title="If enabled and the deployment URL has a path like '/thing', the path will be passed through to the underlying resource instead of being removed.">
                  <InfoCircleOutlined />
                </Tooltip>
              </Checkbox>
            </FormItem>
          ),
        },
      ]}
    />
  );
}

export function createDeployment(values: BaseDeploymentValues) {
  return createDeploymentApiCall({
    body: {
      name: values.name,
      url: values.url,
      externalSource: values.externalSource,
      externalSourceType: values.externalSourceType as "Github" | undefined,
      preserveExternalPath: values.preserveExternalPath,
    },
  });
}
