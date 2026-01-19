import { useState } from "react";
import { useNavigate } from "react-router";
import { Button, Checkbox, Flex, Form, Tooltip, Tree, Typography, Upload, message } from "antd";
import type { TreeDataNode } from "antd";
import { InfoCircleOutlined, UploadOutlined } from "@ant-design/icons";
import { deployFiles } from "~/api-calls/generated/golfComponents";
import { compressFilesToTarGz, getFilesFromUploadList } from "~/utils/fileCompression";
import {
  AdvancedFields,
  CommonDeploymentFields,
  createDeployment,
  type BaseDeploymentValues,
} from "~/components/forms/CommonDeploymentFields";

interface StaticSiteFormValues extends BaseDeploymentValues {
  spaMode?: boolean;
  files?: { originFileObj?: File }[];
}

const FormItem = Form.Item<StaticSiteFormValues>;

function buildFileTree(files: File[], url?: string): TreeDataNode[] {
  const filePaths = files.map((f) => f.webkitRelativePath || f.name).filter((p) => !!p);

  if (filePaths.length === 0) return [];

  const splitPaths = filePaths.map((p) => p.split("/").filter(Boolean));

  // If all files share a single top-level folder (typical for directory uploads),
  // remove it so the tree starts at the site root.
  const firstSegments = splitPaths.map((parts) => parts[0]).filter(Boolean);
  const hasCommonTopLevel =
    firstSegments.length === splitPaths.length &&
    new Set(firstSegments).size === 1 &&
    splitPaths.every((parts) => parts.length > 1);

  const normalizedSplitPaths = hasCommonTopLevel
    ? splitPaths.map((parts) => parts.slice(1))
    : splitPaths;

  const rootChildren: TreeDataNode[] = [];
  const nodeMap = new Map<string, TreeDataNode>();

  for (const parts of normalizedSplitPaths) {
    let currentPath = "";

    for (let i = 0; i < parts.length; i++) {
      const part = parts[i];
      const parentPath = currentPath;
      currentPath = currentPath ? `${currentPath}/${part}` : part;

      if (!nodeMap.has(currentPath)) {
        const isLeaf = i === parts.length - 1;
        const node: TreeDataNode = {
          title: part,
          key: currentPath,
          isLeaf,
          children: isLeaf ? undefined : [],
        };
        nodeMap.set(currentPath, node);

        if (parentPath) {
          const parent = nodeMap.get(parentPath);
          parent?.children?.push(node);
        } else {
          rootChildren.push(node);
        }
      }
    }
  }

  const rootTitle = url?.trim() ? url.trim() : "/";
  return [
    {
      title: rootTitle,
      key: "__site_root__",
      children: rootChildren,
    },
  ];
}

export function StaticSiteForm() {
  const navigate = useNavigate();
  const [form] = Form.useForm<StaticSiteFormValues>();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const uploadFileList = Form.useWatch<StaticSiteFormValues["files"]>("files", form);
  const url = Form.useWatch<StaticSiteFormValues["url"]>("url", form);

  const selectedFiles = getFilesFromUploadList(uploadFileList);
  const selectedFileCount = selectedFiles.length;

  const handleSubmit = async (values: StaticSiteFormValues) => {
    setIsSubmitting(true);
    try {
      await createDeployment(values);

      const files = getFilesFromUploadList(values.files);

      if (files.length > 0) {
        const compressedContents = await compressFilesToTarGz(files);
        await deployFiles({
          body: {
            url: values.url,
            contents: compressedContents,
          },
        });
      }

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

        <Typography.Title level={5}>Deployment Content</Typography.Title>

        <FormItem
          label="Website Content"
          required
          rules={[{ required: true, message: "Please choose a folder to upload" }]}
          name="files"
          valuePropName="fileList"
          getValueFromEvent={(e) => e?.fileList}
        >
          <Upload.Dragger
            multiple
            directory
            accept="*"
            showUploadList={false}
            beforeUpload={() => false}
          >
            <div
              className={
                selectedFileCount > 0
                  ? "relative min-h-64 flex flex-col md:flex-row"
                  : "relative min-h-64 flex"
              }
            >
              <div
                className={
                  selectedFileCount > 0
                    ? "md:w-2/3 flex flex-col items-center justify-center"
                    : "w-full flex flex-col items-center justify-center"
                }
              >
                <p className="ant-upload-drag-icon">
                  <UploadOutlined />
                </p>
                <p className="ant-upload-text">Choose Folder...</p>
                <p className="ant-upload-hint">Select a folder to upload for the site's content.</p>
              </div>

              {selectedFileCount > 0 ? (
                <div
                  className="md:w-1/2 mt-4 md:mt-0 md:pl-4 text-left max-h-56 md:max-h-64 flex flex-col"
                  // prevent clicks on the tree from triggering the file upload dialog:
                  onMouseDown={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                  }}
                  onClick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                  }}
                >
                  <Flex
                    justify="space-between"
                    style={{ backgroundColor: "white", padding: "8px 8px 0 8px" }}
                  >
                    <Typography.Text strong>Selected Files</Typography.Text>
                    <Button
                      size="small"
                      onClick={(e) => {
                        e.preventDefault();
                        e.stopPropagation();
                        form.setFieldValue("files", []);
                      }}
                    >
                      Clear
                    </Button>
                  </Flex>
                  <Tree.DirectoryTree
                    styles={{
                      root: {
                        height: "100%",
                        overflowY: "auto",
                        scrollbarWidth: "thin",
                        padding: 4,
                      },
                    }}
                    showLine
                    defaultExpandAll
                    selectable={false}
                    treeData={buildFileTree(selectedFiles, url)}
                  />
                </div>
              ) : null}
            </div>
          </Upload.Dragger>
        </FormItem>

        <FormItem name="spaMode" valuePropName="checked">
          <Checkbox>
            SPA Mode{" "}
            <Tooltip title="Use /index.html as a fallback for all requests. Useful for single-page applications.">
              <InfoCircleOutlined />
            </Tooltip>
          </Checkbox>
        </FormItem>
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
