package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func runCommand(command string) error {
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stderr = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

const expectedURLStart = "https://go.dev/dl/"

func main() {
	goroot := runtime.GOROOT()
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("You must specify just one link for download")
		return
	}
	link := args[0]
	if !strings.HasPrefix(link, expectedURLStart) {
		fmt.Printf("Download link should start with %q\n", expectedURLStart)
		return
	}
	if len(link) <= len(expectedURLStart) {
		fmt.Printf("Link length should be larger than %q\n", expectedURLStart)
		return
	}
	dlfilename := filepath.Base(link)
	if filepath.Ext(dlfilename) != ".gz" {
		fmt.Printf("Expected download filename extension is %q\n", "gz")
		return
	}

	// 1. Download the version
	err := runCommand(fmt.Sprintf("wget %s", link))
	if err != nil {
		fmt.Println(err)
		return
	}

	// 2. Remove GOROOT folder
	err = runCommand(fmt.Sprintf("sudo rm -rf %s", goroot))
	if err != nil {
		fmt.Println(err)
		return
	}

	// 3. Install the downloaded version
	err = runCommand(fmt.Sprintf("sudo tar -C %s -xzf %s", filepath.Dir(goroot), dlfilename))
	if err != nil {
		fmt.Println(err)
		return
	}

	// 4. Remove the zip installation file
	err = runCommand(fmt.Sprintf("rm %s", dlfilename))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("👍 DONE!")
}
