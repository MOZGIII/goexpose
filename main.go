package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

var gopath = flag.String("gopath", getEnvGoPath(), "GOPATH to work with")
var allowRelative = flag.Bool("allow-relative", false, "does not try to use absolute path to code")
var noProjectNameGuess = flag.Bool("no-guess", false, "disables guessing of project names")

func getEnvGoPath() string {
	value := os.Getenv("GOEXPOSEPATH")
	if value != "" {
		return value
	}
	return os.Getenv("GOPATH")
}

func getLastPath(pathsList string) string {
	split := filepath.SplitList(pathsList)
	return split[len(split)-1]
}

func guessProjectName(codeRootPath string) string {
	codeRootPath = path.Clean(codeRootPath) // need clean path here
	projectName := path.Base(codeRootPath)

	// Passed src dir, that's most probably the repo structure, try looking one dir above
	if projectName == "src" {
		fmt.Println(path.Dir(codeRootPath))
		projectName = path.Base(path.Dir(codeRootPath))
	}
	return projectName
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s /path/to/code [project-name]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "By default environment variable GOPATH is used to determine Go environment.")
		fmt.Fprintln(os.Stderr, "If GOPATH passed consists of multiple parts, the last one is actually used.")
		fmt.Fprintln(os.Stderr, "You can also set GOEXPOSEPATH instead of GOPATH for use with this command.")
	}
	flag.Parse()

	if *gopath == "" {
		fmt.Fprintf(os.Stderr, "Error: passed empty GOPATH!\nYou must pass valid GOPATH environment variable or use -gopath flag to specify Go environment to work with.\n")
		os.Exit(2)
		return
	}

	lastGoPath := filepath.ToSlash(getLastPath(*gopath))

	argc := flag.NArg()
	if argc < 1 || argc > 2 {
		flag.Usage()
		os.Exit(2)
		return
	}

	codeRootPath := filepath.Clean(flag.Arg(0))

	if !*allowRelative {
		var err error
		codeRootPath, err = filepath.Abs(codeRootPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(2)
			return
		}
	}

	codeRootPath = filepath.ToSlash(codeRootPath)

	stat, err := os.Stat(codeRootPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: path \"%s\" does not exist!\n", codeRootPath)
		} else {
			fmt.Fprintf(os.Stderr, "Error: path to code: \"%s\"\n", err.Error())
		}
		os.Exit(2)
		return
	}

	if !stat.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: \"%s\" is not a directory!\n", codeRootPath)
		os.Exit(2)
		return
	}

	projectName := flag.Arg(1)
	if projectName == "" {
		if !*noProjectNameGuess {
			projectName = guessProjectName(codeRootPath)
		}
	}

	if projectName == "." || projectName == "" {
		fmt.Fprintf(os.Stderr, "Error: specify correct project name explicitly!\n")
		os.Exit(2)
		return
	}

	if filepath.IsAbs(projectName) {
		fmt.Fprintf(os.Stderr, "Error: project name (\"%s\") must not be an absolute path!\n", projectName)
		os.Exit(2)
		return
	}

	fmt.Printf("Exposing \"%s\" as \"%s\" at \"%s\"\n", codeRootPath, projectName, lastGoPath)

	if err := os.Symlink(codeRootPath, filepath.Join(lastGoPath, "src", projectName)); err != nil {
		if os.IsExist(err) {
			fmt.Fprintf(os.Stderr, "Error: \"%s\" already exists in GOPATH.\n", projectName)
			os.Exit(1)
			return
		}

		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
		return
	}
}
