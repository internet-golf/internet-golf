package utils

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// returns two random strings; the first is expected to be a token and the
// second is supposed to be a random ID for the token. (the point of that is
// that the ID can be stored in plaintext and used to look up the token later,
// while the token itself will be hashed)
func GetRandomToken() (string, string) {
	id := make([]byte, 4)
	b := make([]byte, 16)
	rand.Read(b)
	rand.Read(id)
	return fmt.Sprintf("%x", b), fmt.Sprintf("%x", id)
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
func HashStream(stream io.ReadSeeker) (string, error) {
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
func ExtractTarGz(gzipStream io.ReadSeeker, baseOutDir string, trimLeadingDirs bool) error {
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
