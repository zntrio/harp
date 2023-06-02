// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package golang

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	exec "golang.org/x/sys/execabs"
)

// PathSeparatorString models the os.PathSeparator as a string.
var PathSeparatorString = string(os.PathSeparator)

// AllPackagesPath denotes all Go packages in a project.
var AllPackagesPath = "." + PathSeparatorString + "..."

// AllCommandsPath denotes all Go application packages in this project.
var AllCommandsPath = strings.Join([]string{".", "cmd", "..."}, PathSeparatorString)

// GoListSourceFilesTemplate provides a standard Go template for querying
// a project's Go source file paths.
var GoListSourceFilesTemplate = "{{$p := .}}{{range $f := .GoFiles}}{{$p.Dir}}/{{$f}}\n{{end}}"

// GoListTestFilesTemplate provides a standard Go template for querying
// a project's Go test file paths.
var GoListTestFilesTemplate = "{{$p := .}}{{range $f := .XTestGoFiles}}{{$p.Dir}}/{{$f}}\n{{end}}"

// CollectedGoFiles represents source and test Go files in a project.
// Populdated with CollectGoFiles().
var CollectedGoFiles = make(map[string]bool)

// CollectedGoSourceFiles represents the set of Go source files in a project.
// Populated with CollectGoFiles().
var CollectedGoSourceFiles = make(map[string]bool)

// CollectedGoTestFiles represents the set of Go test files in a project.
// Populdated with CollectGoFiles().
var CollectedGoTestFiles = make(map[string]bool)

// CollectGoFiles populates CollectedGoFiles, CollectedGoSourceFiles, and CollectedGoTestFiles.
//
// Vendored files are ignored.
func CollectGoFiles() error {
	var sourceOut bytes.Buffer
	var testOut bytes.Buffer

	//nolint:gosec // controlled input
	cmdSource := exec.Command(
		"go",
		"list",
		"-f",
		GoListSourceFilesTemplate,
		AllPackagesPath,
	)
	cmdSource.Stdout = &sourceOut
	cmdSource.Stderr = os.Stderr

	if err := cmdSource.Run(); err != nil {
		return err
	}

	scannerSource := bufio.NewScanner(&sourceOut)

	for scannerSource.Scan() {
		pth := scannerSource.Text()

		CollectedGoFiles[pth] = true
		CollectedGoSourceFiles[pth] = true
	}

	//nolint:gosec // controlled input
	cmdTest := exec.Command(
		"go",
		"list",
		"-f",
		GoListTestFilesTemplate,
		AllPackagesPath,
	)
	cmdTest.Stdout = &testOut
	cmdTest.Stderr = os.Stderr

	if err := cmdTest.Run(); err != nil {
		return err
	}

	scannerTest := bufio.NewScanner(&testOut)

	for scannerTest.Scan() {
		pth := scannerTest.Text()

		CollectedGoFiles[pth] = true
		CollectedGoTestFiles[pth] = true
	}

	return nil
}
