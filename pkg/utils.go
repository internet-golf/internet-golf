package deppy

import (
	"archive/tar"
	"compress/gzip"
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
		nonDefaultDir = path.Join(homeDir, ".deppy")
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

// from https://stackoverflow.com/a/57640231/3962267
func extractTarGz(gzipStream io.Reader, baseOutDir string) error {
	os.MkdirAll(baseOutDir, 0750)

	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return fmt.Errorf("ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)

	var filePaths []string

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		filePaths = append(filePaths, header.Name)
	}

	fmt.Printf("file paths: %v\n", filePaths)

	longestCommonPrefix := getLongestCommonPrefix(filePaths)
	fmt.Printf("longest common prefix: %v\n", longestCommonPrefix)

	lastSlash := strings.LastIndex(longestCommonPrefix, "/")
	if lastSlash != -1 {
		longestCommonPrefix = longestCommonPrefix[0:lastSlash]
	}

	tarReader = tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		// TODO: if there is only one folder at the root of the tar file, get
		// rid of it and treat it as an ignored prefix in all the other paths

		fmt.Printf("file name: %s\n", header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(path.Join(baseOutDir, header.Name), 0755); err != nil {
				return fmt.Errorf("ExtractTarGz: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			if !filepath.IsLocal(header.Name) {
				return fmt.Errorf("ExtractTarGz: File rejected: %s is not a local file path", header.Name)
			}
			outFile, err := os.Create(path.Join(baseOutDir, header.Name))
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
