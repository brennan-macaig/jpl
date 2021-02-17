package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/brennan-macaig/jpl"
	"log"
	"os"
	"strings"
	"time"
)

const (
	app         = "jpl"
	versionFlag = "version"
	buildFlag   = "b"
	testFlag    = "t"
	fileFlag    = "f"
	helpFlag    = "h"
	initFlag    = "i"
	verifyFile  = "verify"
	usageFormat = app + ` (Jet Propulsion Laboratory) - %s

A simple, easy to use build tool.

usage: ` + app + ` [options]

` + app + ` tries to make good assumptions about what you want it to do.
By default, when invoked with no arguments, it will try to find a build file 
(build.yml), and execute the items in the build array one by one, in the order 
they are listed. If one of the items returns non-zero, it will gracefully exit.

options:
%s
`
	defaultYaml = `
# build.yml
# JPL Build Configuration file
version: 1.0

# Configuration Settings
config:
  passOSEnv: false

# Variables
variables:
  VERSION: 1.0.0

# Tasks
build:
  - module: execute
    commands:
      - echo Build Command

test:
  - module: execute
    commands:
      - echo Test Command

`
)

var (
	version string
)

func main() {
	versionPtr := flag.Bool(versionFlag, false, "Print the version and exit")
	buildPtr := flag.Bool(buildFlag, true, "Run commands for a build")
	testPtr := flag.Bool(testFlag, false, "Run commands for a test, if both build and test are given, only test will run")
	filePtr := flag.String(fileFlag, "build.yml", "Specify build file to use")
	initPtr := flag.Bool(initFlag, false, "Create an empty build file and exit")
	verifyPtr := flag.Bool(verifyFile, false, "Verify the contents of the build file and exit")
	help := flag.Bool(helpFlag, false, "Display this help message and exit")

	flag.Parse()

	if *versionPtr {
		fmt.Println(version)
		os.Exit(1)
	}

	if *help {
		_, _ = fmt.Fprintf(os.Stderr, usageFormat, version, allDefaultsToString(flag.CommandLine))
		os.Exit(1)
	}

	if *initPtr {
		// Write the build file
	}

	if *verifyPtr {
		err := jpl.VerifyFile(*filePtr)
		if err != nil {
			fmt.Printf("could not verify build file: %s\n", err.Error())
			os.Exit(-1)
		}
		os.Exit(1)
	}

	// Get buildfile
	bf, _, err := jpl.ReadBuildFile(*filePtr)
	if err != nil {
		log.Fatalf("cannot process buildfile - %s", err.Error())
	}

	start := time.Now()
	if *testPtr {
		*buildPtr = false
		// Run test
		err = jpl.RunModules(bf.Test)
		if err != nil {
			log.Fatalf("Error: %s", err.Error())
		}
	}

	if *buildPtr {
		// Run build
		err = jpl.RunModules(bf.Build)
		if err != nil {
			log.Fatalf("Error: %s", err.Error())
		}
	}
	fmt.Printf("Done - completed in %s\n", time.Since(start).String())
}

func allDefaultsToString(set *flag.FlagSet) string {
	orig := set.Output()
	out := bytes.NewBuffer(nil)
	set.SetOutput(out)
	set.PrintDefaults()
	set.SetOutput(orig)
	return strings.TrimRight(out.String(), "\n")
}
