package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/mholt/archives"
	"github.com/spf13/cobra"
)

func deployContentCommand() *cobra.Command {
	var files string

	deployContent := cobra.Command{
		Use:     "deploy-content [deployment-name]",
		Example: "deploy-content thing.net --files ./dist",
		Short:   "Deploys content",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.TODO()
			fileTree, err := archives.FilesFromDisk(ctx, nil, map[string]string{
				files: "",
			})
			if err != nil {
				panic(err.Error())
			}
			format := archives.CompressedArchive{
				Compression: archives.Gz{},
				Archival:    archives.Tar{},
			}
			tempFile, tempFileErr := os.CreateTemp("", "files-to-deploy")
			if tempFileErr != nil {
				panic(tempFileErr.Error())
			}
			defer os.Remove(tempFile.Name())

			archiveErr := format.Archive(ctx, tempFile, fileTree)
			if archiveErr != nil {
				panic(archiveErr.Error())
			}

			tempFile.Seek(0, 0)

			client := createClient(args[0])

			body, resp, respError := client.
				DefaultAPI.PutDeployFiles(context.TODO()).
				Url(args[0]).
				Contents(tempFile).
				Execute()

			// TODO: handle responses uniformly across commands
			if respError != nil {
				panic(respError.Error())
			}
			if resp.StatusCode != 200 {
				body, _ := io.ReadAll(resp.Body)
				panic("[error from server]: " + string(body))
			}
			if body == nil || !body.Success {
				panic("Did not get success status back from server. Request was to " + resp.Request.URL.String())
			}
			fmt.Println(body.Message)
		},
	}

	deployContent.Flags().StringVar(
		&files, "files", "",
		"Supply a path to a directory with the content you wish to deploy.",
	)

	return &deployContent
}
