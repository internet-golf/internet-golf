package content

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
	"github.com/toBeOfUse/internet-golf/pkg/db"
	"github.com/toBeOfUse/internet-golf/pkg/utils"
)

type FileManager struct{ Settings db.StorageSettings }

// receives a tar
func (f FileManager) TarGzToDeploymentFiles(
	stream io.ReadSeeker, contentName string,
	keepLeadingDirectories bool, _preserveExistingFiles bool,
) (string, error) {
	hash, hashErr := hashStream(stream)
	if hashErr != nil {
		return "", fmt.Errorf("could not hash files for %s", contentName)
	}
	outDir := path.Join(
		f.Settings.DataDirectory,
		slug.Make(contentName),
		hash,
	)
	// weirdly, formData.Contents is a seekable stream, which i'm pretty
	// sure means its entire contents must be being kept in memory so that
	// they can be sought back to (unless it falls back to saving them
	// to disk for large files?) this seems like an annoying limitation
	if tarGzError := extractTarGz(
		stream, outDir, !keepLeadingDirectories,
	); tarGzError != nil {
		return "", tarGzError
	}

	// TODO: if PreserveExistingFiles, find the existing deployment for
	// this url and copy its files into the new directory (if they're
	// not the same directory (i.e. if the hashes are unequal) (which
	// will presumably require getting a reference to the existing
	// deployment from the bus - maybe instead of creating deployContent
	// we should be getting a cursor/active record))

	return outDir, nil
}

// turns the contents of a stream into an md5 hash. seeks the stream back to its
// start before and after computing the hash.
func hashStream(stream io.ReadSeeker) (string, error) {
	hashWriter := md5.New()
	stream.Seek(0, 0)
	defer stream.Seek(0, 0)

	written, err := io.Copy(hashWriter, stream)
	if err != nil {
		return "", err
	}
	fmt.Printf("hashed %v bytes\n", written)
	return hex.EncodeToString(hashWriter.Sum(nil)), nil
}

// function that takes a stream containing .tar.gz data and extracts the files
// and folders within to baseOutDir.
//
// if trimLeadingDirs is true, parent directories at the top level that have no
// siblings and that contain every other file in the tarball within them will be
// discarded (e.g. if the files in the .tar.gz are ["dist/index.html",
// "dist/index.js", "dist/favicon.ico"], it will discard the "dist/" and just
// create the files ["index.html", "index.js", "favicon.ico"]). this is
// generally what you want.
//
// the tar file traversal is heavily referenced from
// https://stackoverflow.com/a/57640231/3962267
//
// TODO: would probably be easier with
// https://github.com/mholt/archives?tab=readme-ov-file#extract-archive
func extractTarGz(gzipStream io.ReadSeeker, baseOutDir string, trimLeadingDirs bool) error {
	os.MkdirAll(baseOutDir, 0750)

	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return fmt.Errorf("ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)

	longestCommonPrefix := ""
	if trimLeadingDirs {

		// perform a first pass over the contents of the tar data that just
		// gathers the paths of the files in it and figures out if there's a
		// common leading prefix that can be removed

		var filePaths []string

		for {
			header, err := tarReader.Next()

			if err == io.EOF {
				break
			}

			if header.Typeflag == tar.TypeReg {
				filePaths = append(filePaths, header.Name)
			}
		}

		longestCommonPrefix = utils.GetLongestCommonPrefix(filePaths)

		lastSlash := strings.LastIndex(longestCommonPrefix, "/")
		if lastSlash != -1 {
			longestCommonPrefix = longestCommonPrefix[0 : lastSlash+1]
		}

		gzipStream.Seek(0, 0)
		uncompressedStream.Reset(gzipStream)
		tarReader = tar.NewReader(uncompressedStream)
	}

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		if len(longestCommonPrefix) >= len(header.Name) {
			fmt.Printf("skipping %s\n", header.Name)
			continue
		}

		itemName := strings.TrimPrefix(header.Name, longestCommonPrefix)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path.Join(baseOutDir, itemName), 0755); err != nil {
				return fmt.Errorf("ExtractTarGz: MkdirAll() failed: %s", err.Error())
			}
		case tar.TypeReg:
			if !filepath.IsLocal(header.Name) {
				return fmt.Errorf("ExtractTarGz: File rejected: %s is not a local file path", header.Name)
			}
			outFile, err := os.Create(path.Join(baseOutDir, itemName))
			if err != nil {
				return fmt.Errorf("ExtractTarGz: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("ExtractTarGz: Copy() failed: %s", err.Error())
			}
			outFile.Close()

		default:
			return fmt.Errorf(
				"ExtractTarGz: unknown type: %v in %v",
				header.Typeflag,
				header.Name)
		}
	}

	return nil
}
