package utils

import (
	"crypto/rand"
	"fmt"
	"net"
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

func GetLongestCommonPrefix(strings []string) string {
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
