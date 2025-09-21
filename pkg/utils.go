package internetgolf

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type StorageSettings struct {
	DataDirectory string
}

// pass in an empty string to use the default data directory
func (s *StorageSettings) Init(nonDefaultDataDirectory string) {
	var dataDirectoryError error
	s.DataDirectory, dataDirectoryError = getDataDirectory(nonDefaultDataDirectory)
	if dataDirectoryError != nil {
		panic("Could not create data directory: " + dataDirectoryError.Error())
	}
	fmt.Printf("Initialized data directory to %s\n", s.DataDirectory)
}

// receives a dataDirectoryPath; translates "$HOME" to the user's home
// directory; creates a directory at the path if it doesn't already exist
func getDataDirectory(dataDirectoryPath string) (string, error) {
	if strings.Index(dataDirectoryPath, "$HOME") != -1 {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", errors.New(
				"could not obtain home directory, so data directory could not be " +
					"created using it. please manually configure the data directory",
			)
		}
		// hopefully this replaceAll doesn't have weird consequences -
		// everything still seems to work here on windows
		homeDir = strings.ReplaceAll(homeDir, "\\", "/")
		dataDirectoryPath = strings.ReplaceAll(dataDirectoryPath, "$HOME", homeDir)
	}

	if _, err := os.Lstat(dataDirectoryPath); err != nil {
		fmt.Printf("Creating data directory at %v\n", dataDirectoryPath)
		// TODO: hope that 0750 works for permissions. can Caddy access the result???
		// will this work recursively?
		if os.Mkdir(dataDirectoryPath, 0750) != nil {
			return "", errors.New("could not create data directory at " + dataDirectoryPath)
		}
	}

	return dataDirectoryPath, nil
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

// getFreePort asks the kernel for a free open port that is ready to use.
// https://gist.github.com/sevkin/96bdae9274465b2d09191384f86ef39d
// exported for use in tests :/
func GetFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
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
