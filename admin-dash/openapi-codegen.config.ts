import { generateSchemaTypes, generateFetchers } from "@openapi-codegen/typescript";
import { defineConfig } from "@openapi-codegen/cli";
export default defineConfig({
  golf: {
    from: {
      relativePath: "../golf-openapi.yaml",
      source: "file",
    },
    outputDir: "app/api-calls/generated",
    to: async (context) => {
      const filenamePrefix = "golf";
      const { schemasFiles } = await generateSchemaTypes(context, {
        filenamePrefix,
      });
      await generateFetchers(context, {
        filenamePrefix,
        schemasFiles,
      });
    },
  },
});
