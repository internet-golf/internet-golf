import { useState } from "react";
import { useNavigate } from "react-router";
import { Button, Form, Typography, message } from "antd";
import {
  AdvancedFields,
  CommonDeploymentFields,
  createDeployment,
  type BaseDeploymentValues,
} from "~/components/forms/CommonDeploymentFields";

export function EmptyDeploymentForm() {
  const navigate = useNavigate();
  const [form] = Form.useForm<BaseDeploymentValues>();
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (values: BaseDeploymentValues) => {
    setIsSubmitting(true);
    try {
      await createDeployment(values);

      message.success("Deployment created successfully");
      navigate("/deployments");
    } catch (err) {
      message.error("Failed to create deployment");
      console.error(err);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Form form={form} layout="vertical" onFinish={handleSubmit}>
      <div className="border border-gray-200 p-4 mb-4">
        <Typography.Title level={5}>Deployment Info</Typography.Title>

        <CommonDeploymentFields />
      </div>

      <AdvancedFields />

      <div className="flex justify-end mt-4">
        <Button type="primary" htmlType="submit" loading={isSubmitting}>
          Create
        </Button>
      </div>
    </Form>
  );
}
