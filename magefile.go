//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build

var extension string = ""

func SetExtension() {
	if runtime.GOOS == "windows" {
		extension = ".exe"
	}
}

func Build() {
	mg.Deps(BuildClient, BuildServer)
}

func BuildServer() error {
	mg.Deps(SetExtension, InstallDeps)
	fmt.Println("Building server...")
	cmd := exec.Command("go", "build", "-o", "golf-server"+extension, "./cmd")
	return cmd.Run()
}

func GenerateClientSdk() error {
	mg.Deps(InstallDeps)
	fmt.Println("Regenerating client code...")
	openapiCmd := exec.Command("go", "run", "./cmd", "openapi")
	err := openapiCmd.Run()
	if err != nil {
		return err
	}
	generateCmd := exec.Command("go", "generate", "./client-cmd")
	err = generateCmd.Run()
	return err
}

func BuildClient() error {
	mg.Deps(SetExtension, InstallDeps, GenerateClientSdk)
	fmt.Println("Building client...")
	cmd := exec.Command("go", "build", "-o", "golf"+extension, "./client-cmd")
	return cmd.Run()
}

func InstallDeps() error {
	fmt.Println("Installing dependencies...")
	cmd := exec.Command("go", "get", ".")
	return cmd.Run()
}

func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("golf.exe")
	os.RemoveAll("golf-server.exe")
	os.RemoveAll("golf")
	os.RemoveAll("golf-server")
}
