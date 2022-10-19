// Copyright (C) 2017 Michał Matczuk
// Use of this source code is governed by an AGPL-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
)

const usage1 string = `Usage: tunnel [OPTIONS] <command> [command args] [...]
options:
`

const usage2 string = `
Commands:
	tunnel qstart               Quick start defined config  

Examples:
	tunnel qstart -host example.arumiot.com -p 3000

Contributions:
	Written by M. Matczuk (mmatczuk@gmail.com)
	Modified by Arum Innovations Pvt. Ltd.
`

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage1)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, usage2)
	}
}

type options struct {
	config   string
	logLevel int
	version  bool
	command  string
	host     string
	port     int
	protocol string
	args     []string
}

func parseArgs() (*options, error) {
	config := flag.String("config", "tunnel.yml", "Path to tunnel configuration file")
	host := flag.String("host", "hookurl.arumiot.com", "Hook Url")
	port := flag.Int("p", 8000, "Local Port to expose")
	protocol := flag.String("protocol", "http", "Protocol")
	logLevel := flag.Int("log-level", 1, "Level of messages to log, 0-3")
	version := flag.Bool("version", false, "Prints tunnel version")
	flag.Parse()
	opts := &options{
		config:   *config,
		host:     *host,
		port:     *port,
		protocol: *protocol,
		logLevel: *logLevel,
		version:  *version,
		command:  flag.Arg(0),
	}

	if opts.version {
		return opts, nil
	}

	switch opts.command {
	case "":
		flag.Usage()
		os.Exit(2)
	case "id", "list":
		opts.args = flag.Args()[1:]
		if len(opts.args) > 0 {
			return nil, fmt.Errorf("list takes no arguments")
		}
	case "start":
		opts.args = flag.Args()[1:]
		if len(opts.args) == 0 {
			return nil, fmt.Errorf("you must specify at least one tunnel to start")
		}
	case "start-all":
		opts.args = flag.Args()[1:]
		if len(opts.args) > 0 {
			return nil, fmt.Errorf("start-all takes no arguments")
		}
	case "qstart":
		opts.args = flag.Args()[1:]
		if len(opts.args) > 0 {
			return nil, fmt.Errorf("qstart takes no arguments")
		}
	default:
		return nil, fmt.Errorf("unknown command %q", opts.command)
	}
	return opts, nil
}
