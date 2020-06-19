//go:generate goderive .

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/caos/orbos/internal/executables"
)

func main() {

	version := flag.String("version", "none", "Path to the git repositorys path to the file containing orbiters current state")
	commit := flag.String("commit", "none", "Path to the git repositorys path to the file containing orbiters current state")
	githubClientID := flag.String("githubclientid", "none", "ClientID used for OAuth with github as store")
	githubClientSecret := flag.String("githubclientsecret", "none", "ClientSecret used for OAuth with github as store")
	orbctldir := flag.String("orbctl", "", "Build orbctl binaries to this directory")
	debug := flag.Bool("debug", false, "Compile executables with debugging features enabled")

	flag.Parse()

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	_, selfPath, _, _ := runtime.Caller(0)
	cmdPath := filepath.Join(filepath.Dir(selfPath), "..")
	path := curryJoinPath(cmdPath)

	if err := executables.PreBuild(executables.PackableBuilds(executables.Build(
		*debug, *commit, *version, *githubClientID, *githubClientSecret,
		executables.Buildable{MainDir: path("nodeagent"), Env: map[string]string{"GOOS": "linux", "GOARCH": "amd64", "CGO_ENABLED": "0"}},
		executables.Buildable{MainDir: path("health"), Env: map[string]string{"GOOS": "linux", "GOARCH": "amd64", "CGO_ENABLED": "0"}},
	))); err != nil {
		panic(err)
	}

	if *orbctldir == "" {
		return
	}

	// Use all available CPUs from now on
	runtime.GOMAXPROCS(runtime.NumCPU())

	orbctlMain := path("orbctl")
	orbctls := executables.Build(
		*debug, *commit, *version, *githubClientID, *githubClientSecret,
		orbctlBin(orbctlMain, *orbctldir, "darwin", "amd64"),
		orbctlBin(orbctlMain, *orbctldir, "freebsd", "amd64"),
		orbctlBin(orbctlMain, *orbctldir, "linux", "amd64"),
		orbctlBin(orbctlMain, *orbctldir, "openbsd", "amd64"),
		orbctlBin(orbctlMain, *orbctldir, "windows", "amd64"),
	)

	var hasErr bool
	for orbctl := range orbctls {
		if _, err := orbctl(); err != nil {
			hasErr = true
		}
	}
	if hasErr {
		panic("Building orbctl failed")
	}
}

func orbctlBin(mainPath, outPath, goos, goarch string) executables.Buildable {

	arch := "x86_64"
	os := strings.ToUpper(goos[0:1]) + goos[1:]
	switch goos {
	case "freebsd":
		os = "FreeBSD"
	case "openbsd":
		os = "OpenBSD"
	}

	outdir := filepath.Join(outPath, fmt.Sprintf("orbctl-%s-%s", os, arch))
	if goos == "windows" {
		outdir += ".exe"
	}

	return executables.Buildable{
		MainDir: mainPath,
		OutDir:  outdir,
		Env:     map[string]string{"GOOS": goos, "GOARCH": goarch, "CGO_ENABLED": "0"},
	}
}

func curryJoinPath(cmdPath string) func(dir string) string {
	return func(dir string) string {
		return filepath.Join(cmdPath, dir)
	}
}
