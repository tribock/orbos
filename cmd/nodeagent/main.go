package main

import (
	"flag"
	"fmt"
	"github.com/caos/orbos/internal/start"
	"github.com/caos/orbos/mntr"
	"os"
)

var gitCommit string
var version string

func main() {

	defer func() {
		r := recover()
		if r != nil {
			os.Stderr.Write([]byte(fmt.Sprintf("\x1b[0;31m%v\x1b[0m\n", r)))
			os.Exit(1)
		}
	}()

	verbose := flag.Bool("verbose", false, "Print logs for debugging")
	printVersion := flag.Bool("version", false, "Print build information")
	repoURL := flag.String("repourl", "", "Repository URL")
	ignorePorts := flag.String("ignore-ports", "", "Comma separated list of firewall ports that are ignored")
	nodeAgentID := flag.String("id", "", "The managed machines ID")

	flag.Parse()

	if *printVersion {
		fmt.Printf("%s %s\n", version, gitCommit)
		os.Exit(0)
	}

	if *repoURL == "" || *nodeAgentID == "" {
		panic("flags --repourl and --id are required")
	}
	monitor := mntr.Monitor{
		OnInfo:   mntr.LogMessage,
		OnChange: mntr.LogMessage,
		OnError:  mntr.LogError,
	}

	if *verbose {
		monitor = monitor.Verbose()
	}

	monitor.WithFields(map[string]interface{}{
		"version":     version,
		"commit":      gitCommit,
		"verbose":     *verbose,
		"repourl":     *repoURL,
		"nodeAgentID": *nodeAgentID,
	}).Info("Node Agent is starting")

	naconfig := &start.NodeAgentConfig{
		GitCommit:   gitCommit,
		NodeAgentID: *nodeAgentID,
		IgnorePorts: *ignorePorts,
		RepoURL:     *repoURL,
	}

	if err := start.NodeAgent(monitor, naconfig); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}
