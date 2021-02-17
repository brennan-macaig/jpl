package jpl

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"text/template"
)

const (
	BuildFileVersion = "1.0"
)

type Buildfile struct {
	Version string    `yaml:"bfversion"`
	Conf    Config    `yaml:"config,omitempty"`
	Test    []Modules `yaml:"test,omitempty"`
	Build   []Modules `yaml:"build,omitempty"`
}

type Variables struct {
	Vars map[string]string `yaml:"variables,omitempty"`
}

type Modules struct {
	Name     string   `yaml:"module,omitempty"`
	Src      []string `yaml:"src,omitempty"`
	Dest     string   `yaml:"dest,omitempty"`
	Commands []string `yaml:"commands,omitempty"`
}

type Config struct {
	PassEnvironment bool `yaml:"passOsEnv,omitempty"`
}

func ReadBuildFile(path string) (Buildfile, Variables, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return Buildfile{}, Variables{}, fmt.Errorf("unable to read build file - %w", err)
	}

	var v Variables
	err = yaml.Unmarshal(yamlFile, &v)
	if err != nil {
		return Buildfile{}, Variables{}, fmt.Errorf("could not get variables from YAML - %w", err)
	}

	tmpl, err := template.New("buildfile").Parse(string(yamlFile))
	if err != nil {
		return Buildfile{}, Variables{}, fmt.Errorf("could not parse template - %w", err)
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, v.Vars)
	if err != nil {
		return Buildfile{}, Variables{}, fmt.Errorf("unable to apply template - %w", err)
	}
	var bf Buildfile
	err = yaml.Unmarshal(b.Bytes(), &bf)
	if err != nil {
		return Buildfile{}, Variables{}, fmt.Errorf("could not unmarshal YAML - %w", err)
	}
	if bf.Version != BuildFileVersion {
		fmt.Printf("%+v\n", bf)
		return Buildfile{}, Variables{}, fmt.Errorf("build file is using a different version! Expected %s, got %s",
			BuildFileVersion, bf.Version)
	}

	return bf, v, nil
}

func VerifyFile(path string) error {
	_, v, err := ReadBuildFile(path)
	if err != nil {
		return err
	}
	fmt.Printf("JPL verify: file accepted\nVariables available: %s\n", v.Vars)
	return nil
}
