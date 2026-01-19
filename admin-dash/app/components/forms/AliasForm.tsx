import { useCallback, useEffect, useState } from "react";
import { useNavigate } from "react-router";
import {
  Button,
  Flex,
  Form,
  Input,
  Segmented,
  Select,
  Spin,
  Tooltip,
  Typography,
  message,
} from "antd";
import { InfoCircleOutlined } from "@ant-design/icons";
import {
  createAlias,
  getDeployments,
  type GetDeploymentsResponse,
} from "~/api-calls/generated/golfComponents";
import {
  AdvancedFields,
  CommonDeploymentFields,
  createDeployment,
  type BaseDeploymentValues,
} from "~/components/forms/CommonDeploymentFields";

type AliasMode = "alias" | "redirect";
type DeploymentsList = GetDeploymentsResponse["deployments"] | null;

interface AliasFormValues extends BaseDeploymentValues {
  aliasMode: AliasMode;
  aliasedToDeployment?: string;
  aliasedToUrl?: string;
}

const FormItem = Form.Item<AliasFormValues>;

export function AliasForm() {
  const navigate = useNavigate();
  const [form] = Form.useForm<AliasFormValues>();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const aliasMode = Form.useWatch<AliasMode>("aliasMode", form);

  console.log("aliasMode", aliasMode);

  const [existingDeployments, setExistingDeployments] = useState<DeploymentsList>(null);
  const [deploymentsLoading, setDeploymentsLoading] = useState(false);

  const loadExistingDeployments = useCallback(async () => {
    if (existingDeployments !== null) return;
    setDeploymentsLoading(true);
    try {
      const response = await getDeployments();
      setExistingDeployments(response.deployments ?? []);
    } catch {
      message.error("Failed to load existing deployments");
    } finally {
      setDeploymentsLoading(false);
    }
  }, [existingDeployments]);

  // fetching data in a useEffect is not great practice, but i don't want to
  // install a data-fetching library yet, just for this...
  useEffect(() => {
    loadExistingDeployments();
  }, [loadExistingDeployments]);

  const handleSubmit = async (values: AliasFormValues) => {
    setIsSubmitting(true);
    try {
      const mode = values.aliasMode ?? "alias";
      const aliasedTo = mode === "alias" ? values.aliasedToDeployment : values.aliasedToUrl;

      if (!aliasedTo) {
        message.error("Please specify what this alias points to");
        setIsSubmitting(false);
        return;
      }

      await createDeployment(values);

      await createAlias({
        body: {
          // TODO: why is this field name capitalized?
          Url: values.url,
          aliasedTo,
          redirect: mode === "redirect",
        },
      });

      message.success("Alias created successfully");
      navigate("/deployments");
    } catch (err) {
      message.error("Failed to create alias");
      console.error(err);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Form
      form={form}
      layout="vertical"
      onFinish={handleSubmit}
      initialValues={{ aliasMode: "alias" }}
    >
      <div className="border border-gray-200 p-4 mb-4">
        <Typography.Title level={5}>Deployment Info</Typography.Title>

        <CommonDeploymentFields />

        <Typography.Title level={5}>Alias Options</Typography.Title>

        <FormItem
          label={
            <Flex align="center" gap="small">
              <span>Mode</span>
              <Tooltip title='"Alias" makes this deployment look identical to another one; "Redirect" makes this deployment redirect users to another URL.'>
                <InfoCircleOutlined />
              </Tooltip>
            </Flex>
          }
          name="aliasMode"
        >
          <Segmented
            options={[
              { label: "Alias", value: "alias" },
              { label: "Redirect", value: "redirect" },
            ]}
          />
        </FormItem>

        {aliasMode === "alias" ? (
          <FormItem required label="To (existing deployment)" name="aliasedToDeployment">
            <Select
              placeholder="Select a deployment"
              loading={deploymentsLoading}
              onFocus={loadExistingDeployments}
              notFoundContent={deploymentsLoading ? <Spin size="small" /> : "No deployments found"}
              options={existingDeployments?.map((d) => ({
                label: d.name,
                value: d.url,
              }))}
            />
          </FormItem>
        ) : (
          <FormItem required label="To (URL)" name="aliasedToUrl">
            <Input placeholder="https://" />
          </FormItem>
        )}
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
