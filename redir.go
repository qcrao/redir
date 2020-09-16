// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	daemon  = flag.Bool("s", false, "run redir service")
	operate = flag.String("op", "create", "operators, create/update/delete/fetch")
	alias   = flag.String("a", "", "alias for a new link")
	link    = flag.String("l", "", "actual link for the alias, optional for delete/fetch")
)

func main() {
	log.SetPrefix(conf.Log)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Usage = usage
	flag.Parse()

	if daemon == nil {
		flag.Usage()
		return
	}
	if *daemon {
		processServer()
		return
	}
	processCmd()
}

func processServer() {
	s := newServer(context.Background())
	s.registerHandler()
	log.Printf("serving %s\n", conf.Addr)
	if err := http.ListenAndServe(conf.Addr, nil); err != nil {
		log.Printf("ListenAndServe %s: %v\n", conf.Addr, err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `usage: redir [-s] [-op <operator> -a <alias> -l <link>]
options:
`)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `example:
redir -s                  run the redir service
redir -a alias -l link    allocate new short link if possible
redir -op fetch -a alias  fetch alias information
`)
	os.Exit(2)
}

func processCmd() {
	if operate == nil || !op(*operate).valid() {
		flag.Usage()
		return
	}
	switch o := op(*operate); o {
	case opCreate, opUpdate:
		if alias == nil || link == nil || *alias == "" || *link == "" {
			flag.Usage()
			return
		}
	case opDelete, opFetch:
		if alias == nil || *alias == "" {
			flag.Usage()
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan bool, 1)
	go func() {
		shortCmd(ctx, op(*operate), *alias, *link)
		done <- true
	}()

	select {
	case <-ctx.Done():
		log.Fatalf("command timeout!")
		return
	case <-done:
		return
	}
}