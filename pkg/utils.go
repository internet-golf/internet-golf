package internetgolf

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func getDataDirectory(nonDefaultDir string) (string, error) {
	if nonDefaultDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", errors.New(
				"could not obtain home directory, so data directory could not be " +
					"automatically created. please manually configure the data directory",
			)
		}
		nonDefaultDir = path.Join(
			// hopefully this replaceAll doesn't have weird consequences -
			// everything still seems to work here on windows
			strings.ReplaceAll(homeDir, "\\", "/"), ".internetgolf", // TODO: extract to constant
		)
	}

	if _, err := os.Lstat(nonDefaultDir); err != nil {
		fmt.Printf("Creating data directory at %v\n", nonDefaultDir)
		// TODO: hope that 0750 works for permissions. can Caddy access the result???
		if os.Mkdir(nonDefaultDir, 0750) != nil {
			return "", errors.New("could not create data directory at " + nonDefaultDir)
		}
	} else {
		fmt.Printf("Found data directory at %v\n", nonDefaultDir)
	}

	return nonDefaultDir, nil
}

func getLongestCommonPrefix(strings []string) string {
	longestCommonPrefix := strings[0]

	for i := 1; i < len(strings) && len(longestCommonPrefix) > 0; i++ {
		path := []rune(strings[i])
		newLongestCommonPrefix := ""
		for j, letter := range longestCommonPrefix {
			if j > len(path)-1 {
				break
			}
			if path[j] == letter {
				newLongestCommonPrefix = newLongestCommonPrefix + string(letter)
			} else {
				break
			}
		}
		longestCommonPrefix = newLongestCommonPrefix
	}
	return longestCommonPrefix
}

func hashStream(stream io.ReadSeeker) (string, error) {
	hashWriter := md5.New()
	b := make([]byte, 1024)
	stream.Seek(0, 0)

	for true {
		n, err := stream.Read(b)
		if n == 0 {
			break
		}
		hashWriter.Write(b[:n])
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return "", fmt.Errorf("could not hash stream: %v", err)
			}
		}
	}

	// cleanup
	stream.Seek(0, 0)

	return hex.EncodeToString(hashWriter.Sum(nil)), nil
}

// from https://stackoverflow.com/a/57640231/3962267
func extractTarGz(gzipStream io.ReadSeeker, baseOutDir string, trimLeadingDirs bool) error {
	os.MkdirAll(baseOutDir, 0750)

	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return fmt.Errorf("ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)

	longestCommonPrefix := ""

	if trimLeadingDirs {

		var filePaths []string

		for {
			header, err := tarReader.Next()

			if err == io.EOF {
				break
			}

			filePaths = append(filePaths, header.Name)
		}

		fmt.Printf("file paths: %v\n", filePaths)

		longestCommonPrefix = getLongestCommonPrefix(filePaths)

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
