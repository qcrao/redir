// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"gopkg.in/yaml.v3"
)

// build info, assign by compile time or runtime.
var (
	Version   string
	BuildTime string
	GoVersion = runtime.Version()
)

type config struct {
	Host  string `yaml:"host"`
	Addr  string `yaml:"addr"`
	Store string `yaml:"store"`
	Log   string `yaml:"log"`
	S     struct {
		Prefix string `yaml:"prefix"`
	} `yaml:"s"`
	X struct {
		Prefix     string `yaml:"prefix"`
		VCS        string `yaml:"vcs"`
		ImportPath string `yaml:"import_path"`
		RepoPath   string `yaml:"repo_path"`
	} `yaml:"x"`
}

func (c *config) parse() {
	f := os.Getenv("REDIR_CONF")
	d, err := ioutil.ReadFile(f)
	if err != nil {
		// Just try again with default setting.
		d, err = ioutil.ReadFile("./config.yml")
		if err != nil {
			log.Fatalf("cannot read configuration, err: %v\n", err)
		}
	}
	err = yaml.Unmarshal(d, c)
	if err != nil {
		log.Fatalf("cannot parse configuration, err: %v\n", err)
	}
}

var conf config

func init() {
	conf.parse()
}
