package jpl

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunModules(mod []Modules) error {
	l := len(mod)
	log.Printf("Found %d modules to run", l)
	for x, m := range mod {
		// Run each module from the file
		curr := x + 1
		log.Printf("[%d / %d] Executing module '%s'", curr, l, m.Name)
		switch m.Name {
		case "execute":
			err := Execute(m)
			if err != nil {
				log.Fatalf("Error from module %s - %s", m.Name, err.Error())
			}
			break
		case "copy":
			err := Copy(m)
			if err != nil {
				log.Fatalf("Error from module %s - %s", m.Name, err.Error())
			}
			break
		default:
			log.Fatalf("Unsupported module '%s'", m.Name)
		}
		log.Printf("[%d / %d] Finished running module '%s'", curr, l, m.Name)
	}

	return nil
}

func Execute(mod Modules) error {
	l := len(mod.Commands)
	for x, str := range mod.Commands {
		curr := x + 1
		log.Printf("[Execute] (%d / %d) Running command", curr, l)
		args := strings.Fields(str)
		err := runCommand(args[0], args[1:])
		if err != nil {
			return fmt.Errorf("command error - %w", err)
		}
		log.Printf("[Execute] (%d / %d) Command ran without error", curr, l)
	}
	return nil
}

func runCommand(cmd string, args []string) error {
	command := exec.Command(cmd, args...)
	command.Env = os.Environ()
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	return command.Run()
}

func Copy(mod Modules) error {
	l := len(mod.Src)
	for x, str := range mod.Src {
		curr := x + 1
		log.Printf("[Copy] (%d / %d) Copying '%s'", curr, l, str)
		srcFile, err := os.Open(str)
		if err != nil {
			return fmt.Errorf("could not open src file - %w", err)
		}
		defer srcFile.Close()

		var destFile *os.File
		if l > 1 {
			destFile, err = os.Create(fmt.Sprintf("%s/%s", mod.Dest, filepath.Base(str)))
			if err != nil {
				return fmt.Errorf("could not open dest file - %w", err)
			}
			defer destFile.Close()
		} else {
			destFile, err = os.Create(mod.Dest)
			if err != nil {
				return fmt.Errorf("could not open dest file - %w", err)
			}
			defer destFile.Close()
		}

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return err
		}

		err = destFile.Sync()
		if err != nil {
			return err
		}
	}
	return nil
}
