// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scan

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/tools/go/buildutil"
	"golang.org/x/vuln/internal/govulncheck"
	"golang.org/x/vuln/internal/vulncheck"
)

type config struct {
	vulncheck.Config
	patterns []string
	analysis govulncheck.AnalysisType
	db       string
	json     bool
	dir      string
	verbose  bool
	tags     []string
	test     bool
}

func (c *Cmd) parseFlags() (*config, error) {
	cfg := &config{}
	var tagsFlag buildutil.TagsFlag
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	flags.BoolVar(&cfg.json, "json", false, "output JSON")
	flags.BoolVar(&cfg.verbose, "v", false, "print a full call stack for each vulnerability")
	flags.BoolVar(&cfg.test, "test", false, "analyze test files. Only valid for source code.")
	flags.StringVar(&cfg.db, "db", "https://vuln.go.dev", "vulnerability database URL")
	flags.Var(&tagsFlag, "tags", "comma-separated `list` of build tags")
	flags.Usage = func() {
		fmt.Fprint(flags.Output(), `usage:
	govulncheck [flags] package...
	govulncheck [flags] binary

`)
		flags.PrintDefaults()
		fmt.Fprintf(flags.Output(), "\n%s\n", detailsMessage)
	}
	addTestFlags(flags, cfg)
	if err := flags.Parse(c.Args[1:]); err != nil {
		return nil, err
	}
	cfg.patterns = flags.Args()
	if len(cfg.patterns) == 0 {
		flags.Usage()
		return nil, ErrMissingArgPatterns
	}
	cfg.analysis = govulncheck.AnalysisSource
	if len(cfg.patterns) == 1 {
		if isFile(cfg.patterns[0]) {
			cfg.analysis = govulncheck.AnalysisBinary
		}
	}
	cfg.tags = tagsFlag
	return cfg, nil
}

func isFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}
