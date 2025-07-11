//go:build ignore

/*
Author: @bynect 2025
Description: Automatically generates cgo directives for embedding CPython
*/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

func cmdExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func genFlags(file string, cflags string, ldflags string) error {
	const template =
`// Code generated by gen_flags.go DO NOT EDIT

package libgopy2

/*
#cgo CFLAGS: %s
#cgo LDFLAGS: %s
*/
import "C"
`
	code := fmt.Sprintf(template, cflags, ldflags)
	return os.WriteFile(file, []byte(code), 0644)
}

func listVersions() string {
	out, err := exec.Command("pkg-config", "--list-all").Output()
	if err != nil {
		fmt.Printf("Failed to find python versions with pkg-config: %v\n", err)
		return ""
	}

	var ver string
	lines := strings.Split(string(out), "\n")

	fmt.Println("Automatically detected these Python versions:")

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		if strings.HasPrefix(parts[0], "python") {
			if parts[0] == "python3-embed" {
				ver = "3"
			} else if ver == "" {
				ver = strings.TrimPrefix(parts[0], "python-")
				ver = strings.TrimSuffix(ver, "-embed")
			}

			fmt.Print("\t")
			fmt.Println(parts[0])
		}
	}

	fmt.Println("Not all versions may be listed, they may still be found by the tool")
	fmt.Println("Ignore the duplicates with the `-embed` suffix")
	fmt.Println()

	return ver
}

func validVer(ver string) bool {
	if len(ver) > 2 {
		for _, c := range ver[2:] {
			if !unicode.IsDigit(c) {
				return false
			}
		}
		return ver[0] == '3' && ver[1] == '.'
	} else {
		return ver == "3" // special case
	}
	return false
}

func manualFlag(auto bool, err error, label string, flag *string) {
	if !auto {
		fmt.Printf("Enter %s: ", label)
		if _, err := fmt.Scanln(flag); err != nil {
			panic(err)
		}
		return
	}

	if err != nil {
		panic(err)
	}

	fmt.Println("Cannot insert manual flag in --auto mode")
	os.Exit(1)
}

func main() {
	auto := false
	if len(os.Args) > 1 && os.Args[1] == "--auto" {
		auto = true
	}

	defVer := listVersions()
	ver := defVer

	for ok := !auto; ok; ok = !validVer(ver) {
		fmt.Printf("Enter the version number [%s]: ", ver)
		if _, err := fmt.Scanln(&ver); err != nil {
			if err.Error() == "unexpected newline" {
				ver = defVer
			} else {
				panic(err)
			}
		}
		ver = strings.TrimSpace(ver)
	}

	sep := "-"
	if ver == "3" {
		sep = ""
	}
	libname := fmt.Sprintf("python%s%s-embed", sep, ver)
	command := fmt.Sprintf("python%s-config", ver)

	if ver == "" {
		fmt.Println("Valid python version not found")
		if auto {
			fmt.Println("This script was run in --auto mode")
			fmt.Println("Run `go run gen_flags.go` to interactively select a version")
		}
		os.Exit(1)
	}

	fmt.Printf("Selected version %s\n", ver)
	fmt.Printf("Corresponding pkg-config entry: %s\n", libname)
	fmt.Println()

	if !cmdExists(command) {
		fmt.Printf("Command %s not found, falling back to pkg-config\n", command)
		command = ""
	}

	var cflags string
	var ldflags string

	if command != "" || cmdExists("pkg-config") {
		var includes []byte
		var libs []byte
		var err error

		if command != "" {
			includes, err = exec.Command(command, "--includes", "--embed").Output()
			if err != nil {
				fmt.Printf("Command %s failed, falling back to pkg-config\n", command)
			}
		}

		if command == "" || err != nil {
			includes, err = exec.Command("pkg-config", "--cflags", libname).Output()
			if err != nil {
				fmt.Printf("Command pkg-config failed: %v\n", err)
				manualFlag(auto, err, "CFLAGS", &cflags)
			}
		}

		if command != "" {
			if libs, err = exec.Command(command, "--libs", "--embed").Output(); err != nil {
				fmt.Printf("Command %s failed, falling back to pkg-config\n", command)
			}
		}

		if command == "" || err != nil {
			if libs, err = exec.Command("pkg-config", "--libs", libname).Output(); err != nil {
				fmt.Printf("Command pkg-config failed: %v\n", err)
				manualFlag(auto, err, "LDFLAGS", &ldflags)
			}
		}

		cflags = string(includes)
		ldflags = string(libs)
	} else {
		fmt.Println("Command pkg-config not found, insert flags manually")
		manualFlag(auto, nil, "CFLAGS", &cflags)
		manualFlag(auto, nil, "LDFLAGS", &ldflags)
	}

	cflags = strings.TrimSpace(string(cflags))
	ldflags = strings.TrimSpace(string(ldflags))

	fmt.Printf("Selected CFLAGS: %s\n", cflags)
	fmt.Printf("Selected LDFLAGS: %s\n", ldflags)

	const file = "flags.go"
	if err := genFlags(file, string(cflags), string(ldflags)); err != nil {
		fmt.Printf("Failed to generate output file %s, aborting\n", file)
		panic(err);
	}

	fmt.Printf("Successfully generated %s\n", file)
}
