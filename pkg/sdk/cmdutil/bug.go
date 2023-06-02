// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmdutil

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"zntr.io/harp/v2/build/version"
	"zntr.io/harp/v2/pkg/sdk/log"

	exec "golang.org/x/sys/execabs"
)

// BugReport generates a bug report body
//

func BugReport() string {
	var buf bytes.Buffer
	buf.WriteString(bugHeader)
	printGoVersion(&buf)
	buf.WriteString("### Does this issue reproduce with the latest release?\n\n\n")
	printEnvDetails(&buf)
	printAppDetails(&buf)
	buf.WriteString(bugFooter)

	return buf.String()
}

const bugHeader = `<!-- Please answer these questions before submitting your issue. Thanks! -->
`

const bugFooter = `### What did you do?
<!--
If possible, provide a recipe for reproducing the error.
A complete runnable program is good.
-->

### What did you expect to see?
### What did you see instead?
`

func printGoVersion(w io.Writer) {
	fmt.Fprintf(w, "### What version of Go are you using (`go version`)?\n\n")
	fmt.Fprintf(w, "<pre>\n")
	fmt.Fprintf(w, "$ go version\n")
	printCmdOut(w, "", "go", "version")
	fmt.Fprintf(w, "</pre>\n")
	fmt.Fprintf(w, "\n")
}

func printEnvDetails(w io.Writer) {
	fmt.Fprintf(w, "### What operating system and processor architecture are you using (`go env`)?\n\n")
	fmt.Fprintf(w, "<details><summary><code>go env</code> Output</summary><br><pre>\n")
	fmt.Fprintf(w, "$ go env\n")
	printCmdOut(w, "", "go", "env")
	printGoDetails(w)
	printOSDetails(w)
	printCDetails(w)
	fmt.Fprintf(w, "</pre></details>\n\n")
}

func printGoDetails(w io.Writer) {
	printCmdOut(w, "GOROOT/bin/go version: ", filepath.Join(runtime.GOROOT(), "bin", "go"), "version")
	printCmdOut(w, "GOROOT/bin/go tool compile -V: ", filepath.Join(runtime.GOROOT(), "bin", "go"), "tool", "compile", "-V")
}

func printOSDetails(w io.Writer) {
	switch runtime.GOOS {
	case "darwin":
		printCmdOut(w, "uname -v: ", "uname", "-v")
		printCmdOut(w, "", "sw_vers")
	case "linux":
		printCmdOut(w, "uname -sr: ", "uname", "-sr")
		printCmdOut(w, "", "lsb_release", "-a")
		printGlibcVersion(w)
	case "openbsd", "netbsd", "freebsd", "dragonfly":
		printCmdOut(w, "uname -v: ", "uname", "-v")
	case "illumos", "solaris":
		// Be sure to use the OS-supplied uname, in "/usr/bin":
		printCmdOut(w, "uname -srv: ", "/usr/bin/uname", "-srv")
		out, err := os.ReadFile("/etc/release")
		if err == nil {
			fmt.Fprintf(w, "/etc/release: %s\n", out)
		}
	}
}

func printAppDetails(w io.Writer) {
	bi := version.NewInfo()
	fmt.Fprintf(w, "### What version of Secret are you using (`harp version`)?\n\n")
	fmt.Fprintf(w, "<pre>\n")
	fmt.Fprintf(w, "$ harp version\n")
	fmt.Fprintf(w, "%s\n", bi.String())
	fmt.Fprintf(w, "</pre>\n")
	fmt.Fprintf(w, "\n")
}

func printCDetails(w io.Writer) {
	printCmdOut(w, "lldb --version: ", "lldb", "--version")
	cmd := exec.Command("gdb", "--version")
	out, err := cmd.Output()
	if err == nil {
		// There's apparently no combination of command line flags
		// to get gdb to spit out its version without the license and warranty.
		// Print up to the first newline.
		fmt.Fprintf(w, "gdb --version: %s\n", firstLine(out))
	}
}

// printCmdOut prints the output of running the given command.
// It ignores failures; 'go bug' is best effort.
func printCmdOut(w io.Writer, prefix, path string, args ...string) {
	cmd := exec.Command(path, args...)
	out, err := cmd.Output()
	if err != nil {
		return
	}
	fmt.Fprintf(w, "%s%s\n", prefix, bytes.TrimSpace(out))
}

// firstLine returns the first line of a given byte slice.
func firstLine(buf []byte) []byte {
	idx := bytes.IndexByte(buf, '\n')
	if idx > 0 {
		buf = buf[:idx]
	}
	return bytes.TrimSpace(buf)
}

// printGlibcVersion prints information about the glibc version.
// It ignores failures.
func printGlibcVersion(w io.Writer) {
	tempdir := os.TempDir()
	if tempdir == "" {
		return
	}
	src := []byte(`int main() {}`)
	srcfile := filepath.Join(tempdir, "go-bug.c")
	outfile := filepath.Join(tempdir, "go-bug")
	err := os.WriteFile(srcfile, src, 0o600)
	if err != nil {
		return
	}
	defer func() {
		log.CheckErr("unable to remove the go-bug.c file", os.Remove(srcfile))
	}()
	//nolint:gosec // controlled inputs
	cmd := exec.Command("gcc", "-o", outfile, srcfile)
	if _, err = cmd.CombinedOutput(); err != nil {
		return
	}
	defer func() {
		log.CheckErr("unable to remove the go-bug file", os.Remove(outfile))
	}()
	//nolint:gosec // controlled inputs
	cmd = exec.Command("ldd", outfile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	re := regexp.MustCompile(`libc\.so[^ ]* => ([^ ]+)`)
	m := re.FindStringSubmatch(string(out))
	if m == nil {
		return
	}
	//nolint:gosec // controlled input
	cmd = exec.Command(m[1])
	out, err = cmd.Output()
	if err != nil {
		return
	}
	fmt.Fprintf(w, "%s: %s\n", m[1], firstLine(out))

	// print another line (the one containing version string) in case of musl libc
	if idx := bytes.IndexByte(out, '\n'); bytes.Contains(out, []byte("musl")) {
		fmt.Fprintf(w, "%s\n", firstLine(out[idx+1:]))
	}
}
