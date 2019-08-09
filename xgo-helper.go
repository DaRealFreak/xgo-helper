package main

import (
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type configuration struct {
	dist    string
	module  string
	source  string
	targets string
	pkg     string
	image   string
	out     string
}

type cmd struct {
	rootCmd       *cobra.Command
	configuration configuration
}

//
func NewCmd() *cmd {
	mainCmd := &cmd{configuration: configuration{}}
	mainCmd.rootCmd = &cobra.Command{
		Use:   "app",
		Short: "Wrapper for custom xgo container cross compiling local/private repository projects",
		Long:  "Application functioning as a wrapper for the docker container diebaumchen/xgo cross compiling local/private projects",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			mainCmd.runDockerContainer(args[0])
		},
	}
	mainCmd.rootCmd.Flags().StringVar(&mainCmd.configuration.dist, "dest", "", "Destination folder to put binaries in (empty = current)")
	mainCmd.rootCmd.Flags().StringVar(&mainCmd.configuration.module, "module", "", "Module name for local compilation (empty = external git repository)")
	mainCmd.rootCmd.Flags().StringVar(&mainCmd.configuration.source, "source", "", "Repository source (branch/tag/commit hash)")
	mainCmd.rootCmd.Flags().StringVar(&mainCmd.configuration.targets, "targets", "", "Build targets")
	mainCmd.rootCmd.Flags().StringVar(&mainCmd.configuration.pkg, "pkg", "", "Package of main.go")
	mainCmd.rootCmd.Flags().StringVar(&mainCmd.configuration.image, "image", "diebaumchen/xgo", "Docker Image used")
	mainCmd.rootCmd.Flags().StringVar(&mainCmd.configuration.out, "out", "", "Prefix to use for output naming")
	return mainCmd
}

// execute the command function
func main() {
	if err := NewCmd().rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

// run the helper docker container
func (cmd *cmd) runDockerContainer(srcDir string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if cmd.configuration.dist == "" {
		cmd.configuration.dist = dir
	} else {
		// docker for windows requires absolute paths
		cmd.configuration.dist, _ = filepath.Abs(cmd.configuration.dist)
	}

	// only use the custom docker container if it's a local/private repository, use original xgo on external repositories
	if _, err := os.Stat(srcDir); !os.IsNotExist(err) {
		srcDir, _ = filepath.Abs(srcDir)
	} else {
		xgoArgs := os.Args[1:]
		for i, arg := range xgoArgs {
			// drop illegal arguments
			if strings.Contains(arg, "-module") || strings.Contains(arg, "-source") {
				xgoArgs = append(xgoArgs[:i], xgoArgs[i+1:]...)
			}
		}

		err = cmd.run(exec.Command("xgo", xgoArgs...))
		if err != nil {
			// error occurred -> status code 1
			os.Exit(1)
		} else {
			// no error occurred -> status code 0
			os.Exit(0)
		}
	}

	args := []string{
		"run", "--rm",
		"-v", srcDir + ":/src",
		"-v", cmd.configuration.dist + ":/build",
		"-e", "MODULE=" + cmd.configuration.module,
	}

	if cmd.configuration.source != "" {
		args = append(args, "-e", "SOURCE="+cmd.configuration.source)
	}
	if cmd.configuration.targets != "" {
		args = append(args, "-e", "TARGETS="+cmd.configuration.targets)
	}
	if cmd.configuration.pkg != "" {
		args = append(args, "-e", "PACKAGE="+cmd.configuration.pkg)
	}
	if cmd.configuration.out != "" {
		args = append(args, "-e", "OUT="+cmd.configuration.out)
	}

	args = append(args, cmd.configuration.image)
	_ = cmd.run(exec.Command("docker", args...))
}

// run command while passing stdout/stderr to the OS stdout/stderr
func (cmd *cmd) run(command *exec.Cmd) error {
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}
