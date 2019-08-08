package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//ToDo: clean up everything

var (
	dist    = flag.String("dist", "", "Destination folder to put binaries in (empty = current)")
	module  = flag.String("module", "", "Module name for local compilation (empty = external git repository)")
	source  = flag.String("source", "", "Repository source (branch/tag/commit hash)")
	targets = flag.String("targets", "", "Build targets")
	pkg     = flag.String("pkg", "", "Package of main.go")
)

func main() {
	flag.Parse()
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if *dist == "" {
		dist = &dir
	} else {
		// docker for windows requires absolute paths
		dir, _ = filepath.Abs(*dist)
		dist = &dir
	}

	srcDir := flag.Args()[0]
	if !strings.HasPrefix(srcDir, "github.com") {
		srcDir, _ = filepath.Abs(srcDir)
	}

	args := []string{
		"run", "--rm",
		"-v", *dist + ":/build",
		"-v", srcDir + ":/src",
		"-e", "MODULE=" + *module,
		"-e", "SOURCE=" + *source,
		"-e", "TARGETS=" + *targets,
		"-e", "PACKAGE=" + *pkg,
	}
	args = append(args, "diebaumchen/xgo")
	_ = run(exec.Command("docker", args...))
}

func run(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
